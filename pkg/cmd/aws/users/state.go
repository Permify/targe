package users

import (
	tea "github.com/charmbracelet/bubbletea"
)

// State represents the users flow state.
type State struct {
	user         *User
	operation    *Operation
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

// GetOperation retrieves the operation from the state.
func (s *State) GetOperation() *Operation {
	return s.operation
}

// GetPolicyOption retrieves the policy option from the state.
func (s *State) GetPolicyOption() *CustomPolicyOption {
	return s.policyOption
}

// GetService retrieves the service from the state.
func (s *State) GetService() *Service {
	return s.service
}

// GetResource retrieves the resource from the state.
func (s *State) GetResource() *Resource {
	return s.resource
}

// GetPolicy retrieves the policy from the state.
func (s *State) GetPolicy() *Policy {
	return s.policy
}

// Setters

// SetUser updates the user in the state.
func (s *State) SetUser(user *User) {
	s.user = user
}

// SetOperation updates the action in the state.
func (s *State) SetOperation(operation *Operation) {
	s.operation = operation
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

// Constants representing user operations and their slugs
const (
	AttachPolicySlug       = "attach_policy"
	DetachPolicySlug       = "detach_policy"
	AddToGroupSlug         = "add_to_group"
	RemoveFromGroupSlug    = "remove_from_group"
	AttachCustomPolicySlug = "attach_custom_policy"
)

// ReachableOperations Predefined list of actions with their names and descriptions
var ReachableOperations = map[string]Operation{
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

// Constants representing custom policy options and their slugs
const (
	WithoutResourceSlug = "without_resource"
	WithResourceSlug    = "with_resource"
)

// ReachableCustomPolicyOptions Predefined list of custom policy options with their names and descriptions
var ReachableCustomPolicyOptions = map[string]CustomPolicyOption{
	WithoutResourceSlug: {
		Name: "Without Resource (without_resource)",
		Desc: "Applies globally without a resource.",
	},
	WithResourceSlug: {
		Name: "With Resource (with_resource)",
		Desc: "Scoped to a specific resource.",
	},
}

// FindFlow determines the next step based on the current state.
func (s *State) Next() tea.Model {
	// Handle case where user is not defined
	if s.user == nil {
		return UserList(s)
	}

	// Handle case where action is not defined
	if s.operation == nil {
		return OperationList(s)
	}

	// Handle specific action: AttachCustomPolicySlug
	if s.operation.Id == AttachCustomPolicySlug {

		if s.policy != nil {
			return Result(s)
		}

		// Handle case where a policy option is selected
		if s.policyOption != nil {
			switch s.policyOption.Id {
			case WithoutResourceSlug:
				return CreatePolicy(s)

			case WithResourceSlug:
				// Handle case where resource is defined
				if s.resource != nil {
					return CreatePolicy(s)
				}

				// Handle case where service is defined
				if s.service != nil {
					return ResourceList(s)
				}
				// If service is not defined
				return ServiceList(s)
			}
		} else {
			// If no policy option is selected
			return CustomPolicyOptionList(s)
		}
	}

	// Handle case where no policy is selected
	if s.policy == nil {
		return PolicyList(s)
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
