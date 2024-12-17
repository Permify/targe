package groups

import (
	tea "github.com/charmbracelet/bubbletea"

	internalaws "github.com/Permify/kivo/internal/aws"
)

// State represents the groups flow state.
type State struct {
	group        *Group
	operation    *Operation
	policyOption *CustomPolicyOption
	service      *Service
	resource     *Resource
	policy       *Policy
}

// Getters

// GetGroup retrieves the group from the state.
func (s *State) GetGroup() *Group {
	return s.group
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

// SetGroup updates the group in the state.
func (s *State) SetGroup(group *Group) {
	s.group = group
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

// Constants representing group actions and their slugs
const (
	AttachPolicySlug       = "attach_policy"
	DetachPolicySlug       = "detach_policy"
	AttachCustomPolicySlug = "attach_custom_policy"
)

// ReachableOperations Predefined list of actions with their names and descriptions
var ReachableOperations = map[string]Operation{
	AttachPolicySlug: {
		Name: "Attach Policy (attach_policy)",
		Desc: "Assign a policy to the group.",
	},
	DetachPolicySlug: {
		Name: "Detach Policy (detach_policy)",
		Desc: "Remove a policy from the group.",
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
func (s *State) Next(api *internalaws.Api) tea.Model {
	// Handle case where group is not defined
	if s.group == nil {
		return GroupList(api, s)
	}

	// Handle case where action is not defined
	if s.operation == nil {
		return OperationList(api, s)
	}

	// Handle specific action: AttachCustomPolicySlug
	if s.operation.Id == AttachCustomPolicySlug {

		if s.policy != nil {
			return Result(api, s)
		}

		// Handle case where a policy option is selected
		if s.policyOption != nil {
			switch s.policyOption.Id {
			case WithoutResourceSlug:
				return CreatePolicy(api, s)

			case WithResourceSlug:
				// Handle case where resource is defined
				if s.resource != nil {
					return CreatePolicy(api, s)
				}

				// Handle case where service is defined
				if s.service != nil {
					return ResourceList(api, s)
				}
				// If service is not defined
				return ServiceList(api, s)
			}
		} else {
			// If no policy option is selected
			return CustomPolicyOptionList(api, s)
		}
	}

	// Handle case where no policy is selected
	if s.policy == nil {
		return PolicyList(api, s)
	}

	// Default fallback
	return Result(api, s)
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
