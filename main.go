package main

import (
	"github.com/gen2brain/beeep"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"log"
	"time"
)

var utcZone, _ = time.LoadLocation("UTC")

func main() {
	nextBreak := time.Now().In(utcZone).Add(time.Minute * 20)

	if err := ui.Init(); err != nil {
		log.Fatalf("Could not intialize the terminal window: %v", err)
	}

	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	renderUi := func(paragraph *widgets.Paragraph) *widgets.Paragraph {
		ui.Clear()

		grid.Set(
			ui.NewRow(1.0,
				ui.NewCol(1.0, paragraph),
			),
		)

		ui.Render(grid)

		return paragraph
	}

	renderUi(infoPanel(nextBreak.Format(time.RFC1123)))

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			if time.Now().In(utcZone).After(nextBreak) {
				err := beeep.Alert(
					"Eye Break",
					"Take 20 seconds to look at something 20 feet away",
					"assets/warning.png")

				if err != nil {
					panic(err)
				}

				nextBreak = nextBreak.Add(time.Minute * 20)

				renderUi(infoPanel(nextBreak.Format(time.RFC1123)))
				continue
			}
		}
		ui.Clear()
		ui.Render(grid)
	}
}

func infoPanel(nextUpdateTime string) *widgets.Paragraph {
	displayLine := "Next eye break: " + nextUpdateTime

	p := widgets.NewParagraph()
	p.Title = "Twenty Twenty Twenty"
	p.TextStyle.Fg = ui.ColorGreen
	p.Text = displayLine
	p.TextStyle.Fg = ui.ColorCyan
	p.BorderStyle.Fg = ui.ColorMagenta

	return p
}
