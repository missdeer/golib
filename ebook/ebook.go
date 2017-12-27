package ebook

// IBook interface for variant ebook generators
type IBook interface {
	Info()
	Begin()
	End()
	SetTitle(string)
	AppendContent(string, string, string)
	SetMargins(float64, float64)
	SetPageType(string)
	SetFontSize(int, int)
	SetLineSpacing(float64)
	SetFontFamily(string)
	SetFontFile(string)
}
