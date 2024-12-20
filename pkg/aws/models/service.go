package models

type Service struct {
	Name string
	Desc string
}

func (i Service) Title() string       { return i.Name }
func (i Service) Description() string { return i.Desc }
func (i Service) FilterValue() string { return i.Name }
