package term

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type keypress struct {
	queryRender string
}

type id struct {
	index int
}

// Home : Renders the initial screen in the terminal
func Home() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	Logo := widgets.NewParagraph()
	// Logo.SetRect(40, 0, 400, 9)
	Logo.Border = false
	Logo.TextStyle.Fg = ui.ColorRed
	Logo.Text = `
			________    _________                                               
			\_____  \  /   _____/ ____   _____   ____       ____   ____   ____  
			 /  ____/  \_____  \ /  _ \ /     \_/ __ \     /  _ \ /    \_/ __ \ 
			/       \  /        (  <_> )  Y Y  \  ___/    (  <_> )   |  \  ___/ 
			\_______ \/_______  /\____/|__|_|  /\___  > /\ \____/|___|  /\___  >`

	title := widgets.NewParagraph()
	title.Text = "-- A Half-online Terminal Sync Play App"
	title.Border = false
	title.TextStyle.Fg = ui.ColorBlue

	permitPane := widgets.NewTabPane("Play a song", "Play a video", "Quit")
	permitPane.PaddingTop = 1
	// permitPane.SetRect(40, 10, 60, 13)
	permitPane.Border = true

	renderTabOne := func() {
		switch permitPane.ActiveTabIndex {
		case 0:
			Play("D:\\8.mp3", "D:\\8.mp3", "Audio")
		case 1:
			Play("D:\\TEMP\\aa.mp4", "Time Sync", "Video")
		case 2:
			ui.Clear()
			ui.Close()
			os.Exit(0)
		}
	}

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewRow(0.6, Logo),
			ui.NewRow(0.4, title),
		),
		ui.NewRow(1.0/6, permitPane),
	)

	ui.Render(grid)

	menuEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-menuEvents:
			switch e.ID {
			case "<Escape>":
				ui.Clear()
				ui.Close()
				os.Exit(0)
			case "<Left>":
				permitPane.FocusLeft()
				ui.Clear()
				ui.Render(grid)
				if e.Type == ui.KeyboardEvent && e.ID == "<Enter>" {
					renderTabOne()
				}
			case "<Right>":
				permitPane.FocusRight()
				ui.Clear()
				ui.Render(grid)
				if e.Type == ui.KeyboardEvent && e.ID == "<Enter>" {
					renderTabOne()
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			case "<Enter>":
				renderTabOne()
			}
		case <-ticker:
			ui.Clear()
			ui.Render(grid)
		}
	}
}
