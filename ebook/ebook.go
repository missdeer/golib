package ebook

// IBook interface for variant ebook generators
type IBook interface {
	Info()
	Begin()
	End()
	SetTitle(string)
	AppendContent(string, string, string)
}
