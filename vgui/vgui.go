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
	"strconv"
	"strings"
	"time"
)

var path string
var pause bool = false
var volume int = 0

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
	//endUpdateProgress := make(chan bool)

	err := vlc.Init("--quiet", "--no-xlib")
	assertErr(err)

	// Create a new player.
	player, err := vlc.NewPlayer()
	assertErr(err)
	player.SetVolume(80)

	// Create a Websocket Client
	wsc, err := wsclient.New()
	if err != nil {
		log.Fatal(err)
	}

	statu := make(chan wsclient.Message)
	go wsc.Read(statu)

	//lblTimeUsed = widget.NewLabel("")
	lblVolume := widget.NewLabel("Volume Now: 80")
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 1
	progress.Value = 0
	label3 := widget.NewLabel("Time: ")
	label3.Alignment = fyne.TextAlignCenter
	go updateTime(progress, player, label3)
	//endUpdateProgress <- true

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.VolumeDownIcon(), func() {
			v, err := player.Volume()
			assertErr(err)
			if v <= 100 && v > 0 {
				player.SetVolume(v - 10)
				t := strconv.Itoa(v - 10)
				text := "Volume Now: " + t
				lblVolume.SetText(text)
				lblVolume.Refresh()
			}

		}),
		widget.NewToolbarAction(theme.VolumeUpIcon(), func() {
			v, err := player.Volume()
			assertErr(err)
			if v < 100 && v >= 0 {
				player.SetVolume(v + 10)
				t := strconv.Itoa(v + 10)
				text := "Volume Now: " + t
				lblVolume.SetText(text)
				lblVolume.Refresh()
			}
		}),
		widget.NewToolbarAction(theme.VolumeMuteIcon(), func() {
			v, err := player.Volume()
			assertErr(err)
			if v != 0 {
				player.SetVolume(volume)
				volume = v
				lblVolume.SetText("Volume Now: 0")
				lblVolume.Refresh()
			} else if v == 0 {
				player.SetVolume(volume)
				volume = 0
				t := strconv.Itoa(volume)
				text := "Volume Now: " + t
				lblVolume.SetText(text)
				lblVolume.Refresh()
			}

		}),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			playerReleaseMedia(player)
			if _, err := player.LoadMediaFromPath(path); err != nil {
				log.Printf("Cannot load selected media: %s\n", err)
				return
			}
			player.Play()

			//endUpdateProgress <- false
		}),
		widget.NewToolbarAction(theme.MediaPauseIcon(), func() {
			if !pause {
				pause = true
				player.SetPause(true)
				//endUpdateProgress <- true
			} else if pause {
				pause = false
				player.SetPause(false)
				//endUpdateProgress <- false
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
	label := widget.NewLabel("Video Sync")
	label.Alignment = fyne.TextAlignCenter
	label2 := widget.NewLabel("Play Mp4")
	label2.Alignment = fyne.TextAlignCenter

	browse_file := widget.NewButton("BrowseFile", func() {
		fd := dialog.NewFileOpen(func(uriReadCloser fyne.URIReadCloser, err error) {
			path = change(uriReadCloser.URI().Path())
			label2.Text = uriReadCloser.URI().Name()
			label2.Refresh()
		}, w)
		fd.Resize(fyne.NewSize(900, 500))
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".mp4"}))
		fd.Show()
	})

	c := container.NewVBox(label, lblVolume, browse_file, label2, toolbar, progress, label3)
	w.SetContent(c)
	w.Resize(fyne.NewSize(1000, 600))
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

func updateTime(p *widget.ProgressBar, vp *vlc.Player, label *widget.Label) {
	for {

		t, _ := vp.MediaPosition()
		p.SetValue(float64(t))
		f, _ := vp.MediaTime()
		k, _ := vp.MediaLength()
		te := strconv.Itoa(f / 1000)
		tx := strconv.Itoa(k / 1000)
		text := te + " // " + tx
		label.SetText(text)
		time.Sleep(1 * time.Second)

	}
}
