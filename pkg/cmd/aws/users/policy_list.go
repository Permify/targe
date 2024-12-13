package users

import (
	"context"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	internalaws "github.com/Permify/kivo/internal/aws"
	"github.com/Permify/kivo/internal/requirements/aws"
)

var policiesStyle = lipgloss.NewStyle().Margin(1, 2)

type Policy struct {
	Arn      string
	Name     string
	Document map[string]interface{}
}

func (i Policy) Title() string       { return i.Name }
func (i Policy) Description() string { return i.Arn }
func (i Policy) FilterValue() string { return i.Name }

type PolicyListModel struct {
	state *State
	list  list.Model
	err   error
}

func PolicyList(state *State) PolicyListModel {
	var items []list.Item
	var m PolicyListModel

	p := aws.Policies{}
	policies, err := p.GetPolicies()

	mp := aws.ManagedPolicies{}
	managedPolicies, err := mp.GetPolicies()

	attachedPolicies, err := internalaws.ListAttachedUserPolicies(context.Background(), state.awsConfig, state.user.Name)
	m.err = err

	switch state.operation.Id {
	case AttachPolicySlug:

		for _, policy := range policies {
			if !slices.Contains(attachedPolicies, policy.Name) {
				items = append(items, Policy{
					Name: policy.Name,
					Arn:  policy.Arn,
				})
			}
		}

		for _, policy := range managedPolicies {
			if !slices.Contains(attachedPolicies, policy.Name) {
				items = append(items, Policy{
					Name: policy.Name,
					Arn:  policy.Arn,
				})
			}
		}

	case DetachPolicySlug:
		inlinePolicies, err := internalaws.ListUserInlinePolicies(context.Background(), state.awsConfig, state.user.Name)
		m.err = err

		for _, name := range inlinePolicies {
			items = append(items, Policy{
				Name: name,
				Arn:  "inline",
			})
		}

		for _, policy := range policies {
			if slices.Contains(attachedPolicies, policy.Name) {
				items = append(items, Policy{
					Name: policy.Name,
					Arn:  policy.Arn,
				})
			}
		}

		for _, policy := range managedPolicies {
			if slices.Contains(attachedPolicies, policy.Name) {
				items = append(items, Policy{
					Name: policy.Name,
					Arn:  policy.Arn,
				})
			}
		}
	}

	m.err = err
	m.state = state
	m.list.Title = "Policies"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m PolicyListModel) Init() tea.Cmd {
	return nil
}

func (m PolicyListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			policy := m.list.SelectedItem().(Policy)
			m.state.SetPolicy(&policy)
			return Switch(m.state.Next(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := policiesStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PolicyListModel) View() string {
	return policiesStyle.Render(m.list.View())
}
