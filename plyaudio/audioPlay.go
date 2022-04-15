package plyaudio

import (
	"log"
	"os"
	"time"

	vlc "github.com/adrg/libvlc-go"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

//PlayAudio : 	Function for audio playback
func PlayAudio(path, name string) {

	if err := vlc.Init("--no-video", "--quiet"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	media, err := player.LoadMediaFromPath(path)
	if err != nil {
		log.Fatal(err)
	}
	defer media.Release()

	err = player.Play()
	if err != nil {
		log.Fatal(err)
	}

	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Init(); err != nil {
		log.Printf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	NameBox := widgets.NewParagraph()
	NameBox.Text = "Audio :: " + name
	NameBox.Border = false
	NameBox.TextStyle.Fg = ui.ColorRed

	volumeBar:=widgets.NewGauge()
	volumeBar.BarColor=ui.ColorGreen

	ctrlList := widgets.NewList()
	ctrlList.Title = "CONTROLS"
	ctrlList.TitleStyle.Modifier = ui.ModifierBold
	ctrlList.TitleStyle.Fg = ui.ColorCyan
	ctrlList.Rows = []string{
		"================================================",
		"[Space]	::	Toogle play/pause the music",
		"  [S]		::	Stop playing the current song and exit",
		"  [→]		::	seek forward 10s",
		"  [←]		::	seek backward 10s",
		"  [↑]		::	add volum 5%",
		"  [↓]		::	reduce volum 5%",
		"  [Q]		::	quit the player",
	}
	ctrlList.TextStyle = ui.NewStyle(ui.ColorRed)
	ctrlList.WrapText = true
	ctrlList.Border = true

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(0.1, NameBox),
		ui.NewRow(0.8,
			ui.NewRow(0.2, volumeBar),
			ui.NewRow(0.8, ctrlList),
		),
		// ui.NewRow(0.1),
	)

	ui.Render(grid)
	// ui.Render(NameBox, ctrlList)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for player.WillPlay() {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "p", "<Space>":
				player.TogglePause()
				break
			case "s", "<Escape>":
				player.Stop()
				ui.Clear()
				ui.Close()
				os.Exit(0)
				break
			case "<Up>":
				volume, err := player.Volume()
				if err != nil {
					log.Fatal(err)
				}
				player.SetVolume(volume + 5)
				break
			case "<Down>":
				volume, err := player.Volume()
				if err != nil {
					log.Fatal(err)
				}
				player.SetVolume(volume - 5)
				break
			case "<Left>", "a":
				sec, err := player.MediaTime()
				if err != nil {
					log.Fatal(err)
				}
				player.SetMediaTime(sec - 10000)
				player.Play()
				break
			case "<Right>", "d":
				sec, err := player.MediaTime()
				if err != nil {
					log.Fatal(err)
				}
				player.SetMediaTime(sec + 10000)
			case "q", "<C-c>":
				player.Release()
				ui.Clear()
				ui.Close()
				os.Exit(0)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			volume, _ := player.Volume()
			volumeBar.Percent = volume
			ui.Render(grid)
		}
	}

	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		close(quit)
	}

	eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	<-quit
}
