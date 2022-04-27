package vgui

import (
	"2-some-one-cli/plyvideo"
	"2-some-one-cli/util"
	"2-some-one-cli/wsclient"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	vlc "github.com/adrg/libvlc-go"
	"github.com/flopp/go-findfont"
	"log"
	"os"
	"strings"
	"time"
)

var path string

// 设置环境变量
func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

var volume int = 0

func VguiStart() {
	a := app.NewWithID("2SOMEone")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("2SOMEone")
	vgui(w)
}

func vgui(w fyne.Window) {
	player, err := plyvideo.NewPlayer()

	// Create a Websocket Client
	wsc, err := wsclient.New()
	if err != nil {
		log.Fatal(err)
	}
	statue := make(chan wsclient.Message)
	go wsc.Read(statue)

	lblVolume := widget.NewLabel("Volume Now: 80")
	label := widget.NewLabel("Media Sync")
	label.Alignment = fyne.TextAlignCenter
	label2 := widget.NewLabel("Play Media")
	label2.Alignment = fyne.TextAlignCenter
	label3 := widget.NewLabel("Time: ")
	label3.Alignment = fyne.TextAlignCenter

	browseFile := setBrowseFile(w, label2)
	form := setUrlform(player, label2)
	toolbar := setToolbar(player, lblVolume, w, wsc, statue)
	progress := setProgress()

	go updateTime(progress, player, label3)

	c := container.NewVBox(label, lblVolume, form, browseFile, label2, toolbar, progress, label3)
	w.SetContent(c)
	w.Resize(fyne.NewSize(1000, 600))
	w.ShowAndRun()
	player.Release()
	err = vlc.Release()
	util.AssertErr(err)
}

func setProgress() (progress *widget.ProgressBar) {
	progress = widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 1
	progress.Value = 0
	return
}

func setUrlform(player *plyvideo.VlcPlayer, label2 *widget.Label) (form *widget.Form) {
	urlentry := widget.NewEntry()
	urlentry.SetPlaceHolder("Input Url")
	form = widget.NewForm(&widget.FormItem{Text: "URL", Widget: urlentry})
	form.OnSubmit = func() {
		if util.IsUrl(urlentry.Text) {
			err := player.Load(urlentry.Text)
			util.AssertErr(err)
			label2.Text = urlentry.Text
			label2.Refresh()
		} else {
			urlentry.SetText("Input URL")
		}
	}

	form.OnCancel = func() {
		urlentry.SetText("Input URL")
	}
	return
}

func setToolbar(player *plyvideo.VlcPlayer, lblVolume *widget.Label, w fyne.Window, wsc *wsclient.WsClient, statue chan wsclient.Message) (toolbar *widget.Toolbar) {
	toolbar = widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.VolumeDownIcon(), func() {
			text := player.DownVolume()
			lblVolume.SetText(text)
			lblVolume.Refresh()
		}),
		widget.NewToolbarAction(theme.VolumeUpIcon(), func() {
			text := player.UpVolume()
			lblVolume.SetText(text)
			lblVolume.Refresh()
		}),
		widget.NewToolbarAction(theme.VolumeMuteIcon(), func() {
			text, volumeNew := player.MuteVolume(volume)
			volume = volumeNew
			lblVolume.SetText(text)
			lblVolume.Refresh()
		}),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			if len(path) != 0 {
				err := player.Load(path)
				util.AssertErr(err)
				player.Play()
				util.AssertErr(err)
				w.Resize(fyne.NewSize(500, 250))
			} else {
				return
			}

		}),
		widget.NewToolbarAction(theme.MediaPauseIcon(), func() {
			player.Pause()
		}),
		widget.NewToolbarAction(theme.MediaFastRewindIcon(), func() {
			player.Rewind()
		}),
		widget.NewToolbarAction(theme.MediaFastForwardIcon(), func() {
			player.Forward()
		}),
		widget.NewToolbarAction(theme.UploadIcon(), func() {
			player.UploadTime(wsc)
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			player.SyncTime(statue)
		}),
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			err := player.Close()
			util.AssertErr(err)
			w.Resize(fyne.NewSize(1000, 600))
		}),
		widget.NewToolbarAction(theme.ViewFullScreenIcon(), func() {
			player.FullScreen()
		}),
		widget.NewToolbarSpacer(),
	)
	return
}

func setBrowseFile(w fyne.Window, label2 *widget.Label) (browseFile *widget.Button) {
	browseFile = widget.NewButton("BrowseFile", func() {
		fd := dialog.NewFileOpen(func(uriReadCloser fyne.URIReadCloser, err error) {
			if uriReadCloser == nil {
				log.Println("Cancelled")
				return
			}
			path = util.ConvPath(uriReadCloser.URI().Path())
			label2.Text = uriReadCloser.URI().Name()
			label2.Refresh()
		}, w)
		fd.Resize(fyne.NewSize(900, 500))
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".mp4", ".mkv"}))
		fd.Show()
	})
	return
}

func updateTime(p *widget.ProgressBar, player *plyvideo.VlcPlayer, label *widget.Label) {
	for {
		text, ps := player.TimeNow()
		p.SetValue(ps)
		label.SetText(text)
		time.Sleep(1 * time.Second)
	}
}
