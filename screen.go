package main

import(
	"github.com/rivo/tview"
)

type OutputScreen struct {
	*tview.TextView
	title   string
	content string
}

func NewOutputScreen() *OutputScreen {
	return &OutputScreen{
		TextView: tview.NewTextView(),
		title: "Ping",
	}
}

func (s OutputScreen) UpdateTitle() {
	s.SetTitle(s.title)
}

func (s OutputScreen) UpdateContent() {
	s.SetText(s.content)
	s.ScrollToBeginning()
	app.Draw()
}
