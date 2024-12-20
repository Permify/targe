package models

type Role struct {
	Arn  string
	Name string
}

func (i Role) Title() string       { return i.Name }
func (i Role) Description() string { return i.Arn }
func (i Role) FilterValue() string { return i.Name }
