package main

type Screen struct {
	title   string
	content string
}

func (s Screen) UpdateTitle() {
	output_box.SetTitle(s.title)
}

func (s Screen) UpdateContent() {
	output_box.SetText(s.content)
	output_box.ScrollToBeginning()
	app.Draw()
}
