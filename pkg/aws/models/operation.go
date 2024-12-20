package models

type Operation struct {
	Id   string
	Name string
	Desc string
}

func (i Operation) Title() string       { return i.Name }
func (i Operation) Description() string { return i.Desc }
func (i Operation) FilterValue() string { return i.Name }
