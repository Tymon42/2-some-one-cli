package plyvideo

import (
	"log"
	"os"

	vlc "github.com/adrg/libvlc-go"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

//PlayVideo : 	Function for video playback
func PlayVideo(path, name string) {
	if err := vlc.Init("--quiet"); err != nil {
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

	// media, err := player.LoadMediaFromURL(path)
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
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	NameBox := widgets.NewParagraph()
	NameBox.Text = "Vedio :: " + name
	NameBox.Border = false
	NameBox.TextStyle.Fg = ui.ColorRed
	// NameBox.SetRect(10, 9, 149, 12)

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
		"  [Q]		::	quit the player",
	}
	ctrlList.TextStyle = ui.NewStyle(ui.ColorRed)
	ctrlList.WrapText = true
	// ctrlList.SetRect(20, 17, 110, 25)
	ctrlList.Border = true

	// ui.Render(NameBox, ctrlList)
	grid:=ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/5, NameBox),
		ui.NewRow(2.0/3,
			ui.NewRow(1.0/3, ctrlList),
		),
	)

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for player.WillPlay() {
		e := <-uiEvents
		switch e.ID {
		case "C-p":
			// if player.IsPlaying() {
			// player.SetPause(true)
			player.TogglePause()
			// }
			// TODO: Send Pause sign
			break
		case "<Space>":
			// if player.IsPlaying() == false {
				// player.SetPause(false)
			player.TogglePause()
			// }
			// TODO: Send Play sign
			break
		
		case "C-r":
			
			break
		case "s", "<Escape>":
			player.Stop()
			ui.Clear()
			ui.Close()
			os.Exit(0)
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
