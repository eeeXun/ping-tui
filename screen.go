package main

import(
	"github.com/rivo/tview"
)

type OutputScreen struct {
	*tview.TextView
	Title   string
	Content string
}

func NewOutputScreen() *OutputScreen {
	return &OutputScreen{
		TextView: tview.NewTextView(),
		Title: "Ping",
	}
}

func (s OutputScreen) UpdateTitle() {
	s.SetTitle(s.Title)
}

// Concurrency(app.Draw), do not call in main thread
func (s OutputScreen) RefreshContent() {
	s.SetText(s.Content)
	s.ScrollToBeginning()
	app.Draw()
}

func (s *OutputScreen) ClearContent() {
	s.Content = ""
	s.SetText(s.Content)
}
