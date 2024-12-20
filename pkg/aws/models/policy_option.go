package models

type PolicyOption struct {
	Id   string
	Name string
	Desc string
}

func (i PolicyOption) Title() string       { return i.Name }
func (i PolicyOption) Description() string { return i.Desc }
func (i PolicyOption) FilterValue() string { return i.Name }
