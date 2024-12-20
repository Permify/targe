package models

type Policy struct {
	Arn      string
	Name     string
	Document map[string]interface{}
}

func (i Policy) Title() string       { return i.Name }
func (i Policy) Description() string { return i.Arn }
func (i Policy) FilterValue() string { return i.Name }
