package models

type File struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Content string `json:"content"`
}

func (f *File) SetContent(content string) *File {
	f.Content = content
	return f
}

func (f *File) SetName(name string) *File {
	f.Name = name
	return f
}
