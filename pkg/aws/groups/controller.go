package groups

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Permify/kivo/internal/aws"
	requirements "github.com/Permify/kivo/internal/requirements/aws"
	"github.com/Permify/kivo/pkg/aws/models"
)

type Controller struct {
	api   *aws.Api
	State *State
}

func NewController(api *aws.Api, state *State) *Controller {
	return &Controller{
		api:   api,
		State: state,
	}
}

// FailedMsg represents a failure operation.
type FailedMsg struct {
	Err error
}

type GroupLoadedMsg struct{ List []list.Item }

// LoadGroups loads groups from the AWS API.
func (c *Controller) LoadGroups() tea.Cmd {
	return func() tea.Msg {
		var items []list.Item
		output, err := c.api.ListGroups(context.Background())
		if err != nil {
			return FailedMsg{Err: err}
		}

		for _, group := range output.Groups {
			items = append(items, models.Group{
				Name: *group.GroupName,
				Arn:  *group.Arn,
			})
		}

		return GroupLoadedMsg{List: items}
	}
}

type OperationLoadedMsg struct{ List []list.Item }

// LoadOperations loads operations.
func (c *Controller) LoadOperations() tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{
			models.Operation{Id: AttachPolicySlug, Name: ReachableOperations[AttachPolicySlug].Name, Desc: ReachableOperations[AttachPolicySlug].Desc},
			models.Operation{Id: DetachPolicySlug, Name: ReachableOperations[DetachPolicySlug].Name, Desc: ReachableOperations[DetachPolicySlug].Desc},
			models.Operation{Id: AttachCustomPolicySlug, Name: ReachableOperations[AttachCustomPolicySlug].Name, Desc: ReachableOperations[AttachCustomPolicySlug].Desc},
		}
		return OperationLoadedMsg{List: items}
	}
}

type ServiceLoadedMsg struct{ List []list.Item }

// LoadServices loads services.
func (c *Controller) LoadServices() tea.Cmd {
	return func() tea.Msg {
		t := requirements.Types{}
		services, err := t.GetServices()
		if err != nil {
			return FailedMsg{Err: err}
		}

		var items []list.Item

		for _, service := range services {
			items = append(items, models.Service{
				Name: service.Name,
				Desc: service.Description,
			})
		}
		return ServiceLoadedMsg{List: items}
	}
}

type ResourceLoadedMsg struct{ List []list.Item }

// LoadResources loads resources.
func (c *Controller) LoadResources() tea.Cmd {
	return func() tea.Msg {
		var items []list.Item

		resources, err := c.api.ListResources(c.State.GetService().Name)
		if err != nil {
			return FailedMsg{Err: err}
		}

		for _, resource := range resources {
			items = append(items, models.Resource{
				Name: resource.Name,
				Arn:  resource.Arn,
			})
		}

		return ResourceLoadedMsg{List: items}
	}
}

type PolicyLoadedMsg struct{ List []list.Item }

// LoadPolicies loads policies.
func (c *Controller) LoadPolicies() tea.Cmd {
	return func() tea.Msg {
		var items []list.Item

		policies, err := c.api.ListPolicies(context.Background())
		if err != nil {
			return FailedMsg{Err: err}
		}

		mp := requirements.ManagedPolicies{}
		managedPolicies, err := mp.GetPolicies()
		if err != nil {
			return FailedMsg{Err: err}
		}

		attachedPolicies, err := c.api.ListAttachedGroupPolicies(context.Background(), c.State.group.Name)
		if err != nil {
			return FailedMsg{Err: err}
		}

		switch c.State.operation.Id {
		case AttachPolicySlug:
			for _, policy := range policies.Policies {
				if !slices.Contains(attachedPolicies, *policy.PolicyName) {
					items = append(items, models.Policy{
						Name: *policy.PolicyName,
						Arn:  *policy.Arn,
					})
				}
			}

			for _, policy := range managedPolicies {
				if !slices.Contains(attachedPolicies, policy.Name) {
					items = append(items, models.Policy{
						Name: policy.Name,
						Arn:  policy.Arn,
					})
				}
			}
		case DetachPolicySlug:
			inlinePolicies, err := c.api.ListGroupInlinePolicies(context.Background(), c.State.group.Name)
			if err != nil {
				return FailedMsg{Err: err}
			}

			for _, name := range inlinePolicies {
				items = append(items, models.Policy{
					Name: name,
					Arn:  "inline",
				})
			}

			for _, policy := range policies.Policies {
				if slices.Contains(attachedPolicies, *policy.PolicyName) {
					items = append(items, models.Policy{
						Name: *policy.PolicyName,
						Arn:  *policy.Arn,
					})
				}
			}

			for _, policy := range managedPolicies {
				if slices.Contains(attachedPolicies, policy.Name) {
					items = append(items, models.Policy{
						Name: policy.Name,
						Arn:  policy.Arn,
					})
				}
			}
		}

		return PolicyLoadedMsg{List: items}
	}
}

type PolicyOptionLoadedMsg struct{ List []list.Item }

// LoadPolicyOptions loads operations.
func (c *Controller) LoadPolicyOptions() tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{
			models.PolicyOption{Id: WithoutResourceSlug, Name: ReachablePolicyOptions[WithoutResourceSlug].Name, Desc: ReachablePolicyOptions[WithoutResourceSlug].Desc},
			models.PolicyOption{Id: WithResourceSlug, Name: ReachablePolicyOptions[WithResourceSlug].Name, Desc: ReachablePolicyOptions[WithResourceSlug].Desc},
		}
		return PolicyOptionLoadedMsg{List: items}
	}
}

// Constants representing group actions and their slugs
const (
	AttachPolicySlug       = "attach_policy"
	DetachPolicySlug       = "detach_policy"
	AttachCustomPolicySlug = "attach_custom_policy"
)

// ReachableOperations Predefined list of actions with their names and descriptions
var ReachableOperations = map[string]models.Operation{
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

// ReachablePolicyOptions Predefined list of custom policy options with their names and descriptions
var ReachablePolicyOptions = map[string]models.PolicyOption{
	WithoutResourceSlug: {
		Name: "Without Resource (without_resource)",
		Desc: "Applies globally without a resource.",
	},
	WithResourceSlug: {
		Name: "With Resource (with_resource)",
		Desc: "Scoped to a specific resource.",
	},
}

// Next determines the next step based on the current state.
func (c *Controller) Next() tea.Model {
	// Handle case where group is not defined
	if c.State.group == nil {
		return NewGroupList(c)
	}

	// Handle case where action is not defined
	if c.State.operation == nil {
		return NewOperationList(c)
	}

	// Handle specific action: AttachCustomPolicySlug
	if c.State.operation.Id == AttachCustomPolicySlug {

		if c.State.policy != nil {
			return NewResult(c)
		}

		// Handle case where a policy option is selected
		if c.State.policyOption != nil {
			switch c.State.policyOption.Id {
			case WithoutResourceSlug:
				return NewCreatePolicy(c)

			case WithResourceSlug:
				// Handle case where resource is defined
				if c.State.resource != nil {
					return NewCreatePolicy(c)
				}

				// Handle case where service is defined
				if c.State.service != nil {
					return NewResourceList(c)
				}
				// If service is not defined
				return NewServiceList(c)
			}
		} else {
			// If no policy option is selected
			return NewPolicyOptionList(c)
		}
	}

	// Handle case where no policy is selected
	if c.State.policy == nil {
		return NewPolicyList(c)
	}

	// Default fallback
	return NewResult(c)
}

func (c *Controller) Done() error {
	switch c.State.operation.Id {
	case AttachPolicySlug:
		return c.api.AttachPolicyToGroup(context.Background(), c.State.GetPolicy().Arn, c.State.GetGroup().Name)
	case DetachPolicySlug:
		return c.api.DetachPolicyFromGroup(context.Background(), c.State.GetPolicy().Arn, c.State.GetGroup().Name)
	case AttachCustomPolicySlug:
		jsonBytes, err := json.Marshal(c.State.GetPolicy().Document)
		if err != nil {
			return err
		}
		output, err := c.api.CreatePolicy(context.Background(), c.State.GetPolicy().Name, string(jsonBytes))
		if err != nil {
			return err
		}
		return c.api.AttachPolicyToGroup(context.Background(), *output.Policy.Arn, c.State.GetGroup().Name)
	default:
		return errors.New("operation not supported")
	}

	return nil
}

// Switch handles window size changes and updates the model accordingly.
func Switch(model tea.Model, width, height int) (tea.Model, tea.Cmd) {
	// Always initialize the model
	initCmd := model.Init()

	// Handle window size updates
	if width == 0 && height == 0 {
		return model, initCmd
	}

	updateModel, updateCmd := model.Update(tea.WindowSizeMsg{
		Width:  width,
		Height: height,
	})

	// Combine initialization and update commands
	return updateModel, tea.Batch(initCmd, updateCmd)
}
