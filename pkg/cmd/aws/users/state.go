package users

import (
	tea "github.com/charmbracelet/bubbletea"
)

// State represents the application state.
type State struct {
	user         *User
	action       *UserAction
	policyOption *CustomPolicyOption
	service      *Service
	resource     *Resource
	policy       *Policy
}

// Getters

// GetUser retrieves the user from the state.
func (s *State) GetUser() *User {
	return s.user
}

// GetAction retrieves the action from the state.
func (s *State) GetAction() *UserAction {
	return s.action
}

// GetPolicyOption retrieves the policy option from the state.
func (s *State) GetPolicyOption() *CustomPolicyOption {
	return s.policyOption
}

// GetPolicy retrieves the policy from the state.
func (s *State) GetPolicy() *Policy {
	return s.policy
}

// GetService retrieves the service from the state.
func (s *State) GetService() *Service {
	return s.service
}

// GetResource retrieves the resource from the state.
func (s *State) GetResource() *Resource {
	return s.resource
}

// Setters

// SetUser updates the user in the state.
func (s *State) SetUser(user *User) {
	s.user = user
}

// SetAction updates the action in the state.
func (s *State) SetAction(action *UserAction) {
	s.action = action
}

// SetPolicyOption updates the policy option in the state.
func (s *State) SetPolicyOption(policyOption *CustomPolicyOption) {
	s.policyOption = policyOption
}

// SetService updates the service in the state.
func (s *State) SetService(service *Service) {
	s.service = service
}

// SetResource updates the resource in the state.
func (s *State) SetResource(resource *Resource) {
	s.resource = resource
}

// SetPolicy updates the policy in the state.
func (s *State) SetPolicy(policy *Policy) {
	s.policy = policy
}

// Constants representing user actions and their slugs
const (
	AttachPolicySlug       = "attach_policy"
	DetachPolicySlug       = "detach_policy"
	AddToGroupSlug         = "add_to_group"
	RemoveFromGroupSlug    = "remove_from_group"
	AttachCustomPolicySlug = "attach_custom_policy"
)

// ReachableActions Predefined list of actions with their names and descriptions
var ReachableActions = map[string]UserAction{
	AttachPolicySlug: {
		Name: "Attach Policy (attach_policy)",
		Desc: "Assign a policy to the user.",
	},
	DetachPolicySlug: {
		Name: "Detach Policy (detach_policy)",
		Desc: "Remove a policy from the user.",
	},
	AddToGroupSlug: {
		Name: "Add to Group (add_to_group)",
		Desc: "Include the user in a group.",
	},
	RemoveFromGroupSlug: {
		Name: "Remove from Group (remove_from_group)",
		Desc: "Exclude the user from a group.",
	},
	AttachCustomPolicySlug: {
		Name: "Attach Custom Policy (attach_custom_policy)",
		Desc: "Create and attach a custom policy.",
	},
}

// FindFlow determines the next step based on the current state.
func (s *State) FindFlow() tea.Model {
	// Handle case where user is not defined
	if s.user == nil {
		return Users(s)
	}

	// Handle case where action is not defined
	if s.action == nil {
		return Actions(s)
	}

	// Handle specific action: AttachCustomPolicySlug
	if s.action.Title() == ReachableActions[AttachCustomPolicySlug].Name {
		// Handle case where a policy option is selected
		if s.policyOption != nil {
			switch s.policyOption.Name {
			case "Policy Without Resource (policy_without_resource)":
				return CreatePolicy(s)

			case "Policy With Resource (policy_with_resource)":
				// Handle case where service is defined
				if s.service != nil {
					return Resources(s)
				}
				// If service is not defined
				return Services(s)
			}
		} else {
			// If no policy option is selected
			return CustomPolicyOptions(s)
		}
	}

	// Handle case where no policy is selected
	if s.policy == nil {
		return Policies(s)
	}

	// Default fallback
	return Result(s)
}

// Switch handles window size changes and updates the model accordingly.
func Switch(model tea.Model, width, height int) (tea.Model, tea.Cmd) {
	if width == 0 && height == 0 {
		return model, model.Init()
	}

	return model.Update(tea.WindowSizeMsg{
		Width:  width,
		Height: height,
	})
}
