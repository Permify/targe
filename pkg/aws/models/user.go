package models

type User struct {
	Arn  string
	Name string
}

func (i User) Title() string       { return i.Name }
func (i User) Description() string { return i.Arn }
func (i User) FilterValue() string { return i.Name }
