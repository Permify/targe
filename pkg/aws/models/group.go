package models

type Group struct {
	Arn  string
	Name string
}

func (i Group) Title() string       { return i.Name }
func (i Group) Description() string { return i.Arn }
func (i Group) FilterValue() string { return i.Name }
