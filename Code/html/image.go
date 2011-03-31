package main

type Image struct {
	Width, Height string
	Source, AltText string
}

func NewImage(src, width, height string) (img Image) {
	img.Width = width
	img.Height = height
	img.Source = src
	return
}

func (me Image) String() string {
	return "<img src=\"" + me.Source + "\" alt=\"" + me.AltText + "\"" + Attribute("width", me.Width) + Attribute("height", me.Height) + "/>"
}
