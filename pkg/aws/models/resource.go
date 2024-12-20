package models

type Resource struct {
	Arn  string
	Name string
}

func (i Resource) Title() string       { return i.Name }
func (i Resource) Description() string { return i.Arn }
func (i Resource) FilterValue() string { return i.Arn }
