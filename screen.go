package main

import(
	"github.com/rivo/tview"
)

type OutputScreen struct {
	*tview.TextView
	Title   string
	Text string
}

func NewOutputScreen() *OutputScreen {
	return &OutputScreen{
		TextView: tview.NewTextView(),
		Title: "Ping",
	}
}

func (o *OutputScreen) AddText(s string) {
	if len(o.Text) == 0 {
		o.Text = s
	} else {
		o.Text = s + "\n" + o.Text
	}
}

func (o OutputScreen) UpdateTitle() {
	o.SetTitle(o.Title)
}

// Concurrency(app.Draw), do not call in main thread
func (o OutputScreen) RefreshText() {
	o.SetText(o.Text)
	o.ScrollToBeginning()
	app.Draw()
}

func (o *OutputScreen) ClearText() {
	o.Text = ""
	o.SetText(o.Text)
}
