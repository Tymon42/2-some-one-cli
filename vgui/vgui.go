package vgui

import (
	"2-some-one-cli/wsclient"
	"fmt"
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
	"regexp"
	"strings"
	"time"
)

var path string
var pause bool = false

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

func assertErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func assertConv(ok bool) {
	if !ok {
		log.Panic("invalid widget conversion")
	}
}

func playerReleaseMedia(player *vlc.Player) {
	player.Stop()
	if media, _ := player.Media(); media != nil {
		media.Release()
	}
}

func vgui(w fyne.Window) {
	err := vlc.Init("--quiet", "--no-xlib")
	assertErr(err)

	// Create a new player.
	player, err := vlc.NewPlayer()
	assertErr(err)

	// Create a Websocket Client
	wsc, err := wsclient.New()
	if err != nil {
		log.Fatal(err)
	}

	statu := make(chan wsclient.Message)
	go wsc.Read(statu)

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			playerReleaseMedia(player)

			if _, err := player.LoadMediaFromPath(path); err != nil {
				log.Printf("Cannot load selected media: %s\n", err)
				return
			}
			player.Play()
		}),
		widget.NewToolbarAction(theme.MediaPauseIcon(), func() {
			if !pause {
				pause = true
				player.SetPause(true)
			} else if pause {
				pause = false
				player.SetPause(false)
			}
		}),
		widget.NewToolbarAction(theme.MediaFastRewindIcon(), func() {
			sec, err := player.MediaTime()
			if err != nil {
				log.Fatal(err)
			}
			player.SetMediaTime(sec - 10000)
			player.Play()
		}),
		widget.NewToolbarAction(theme.MediaFastForwardIcon(), func() {
			sec, err := player.MediaTime()
			if err != nil {
				log.Fatal(err)
			}
			player.SetMediaTime(sec + 10000)
		}),
		widget.NewToolbarAction(theme.UploadIcon(), func() {
			nowTime := time.Now().UnixNano() / 1e6
			sec, err := player.MediaTime()
			if err != nil {
				log.Fatal(err)
			}
			var msg wsclient.Message
			msg.MediaTime = sec
			msg.Ts = nowTime
			wsc.Writer(msg)
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			Statu := <-statu
			nowTime := time.Now().UnixNano() / 1e6
			settime := int(nowTime-Statu.Ts) + (Statu.MediaTime)
			player.SetMediaTime(settime)
		}),
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			player.Stop()
		}),
		widget.NewToolbarSpacer(),
	)
	label := widget.NewLabel("Video Mp4")
	label.Alignment = fyne.TextAlignCenter
	label2 := widget.NewLabel("Play Mp4")
	label2.Alignment = fyne.TextAlignCenter
	browse_file := widget.NewButton("BrowseFile", func() {
		fd := dialog.NewFileOpen(func(uriReadCloser fyne.URIReadCloser, err error) {
			path = change(uriReadCloser.URI().Path())
			label2.Text = uriReadCloser.URI().Name()
			label2.Refresh()
		}, w)
		fd.Show()
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".mp4"}))
	})

	c := container.NewVBox(label, browse_file, label2, toolbar)
	w.SetContent(c)
	w.Resize(fyne.NewSize(1280, 720))
	w.ShowAndRun()
	player.Release()
	vlc.Release()
}

func VguiStart() {
	a := app.NewWithID("2SOMEone")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("2SOMEone")
	vgui(w)

}

func change(data string) (pathc string) {
	re3, _ := regexp.Compile("/")
	rep := re3.ReplaceAllStringFunc(data, strings.ToUpper)
	fmt.Println(rep)
	rep2 := re3.ReplaceAllString(data, "\\\\")
	fmt.Println(rep2)

	return rep2
}
