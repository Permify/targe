package users

import (
	"context"
	"errors"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Permify/targe/internal/aws"
	requirements "github.com/Permify/targe/internal/requirements/aws"
	"github.com/Permify/targe/pkg/aws/models"
)

type Controller struct {
	api          *aws.Api
	openAiApiKey string
	State        *State
}

func NewController(api *aws.Api, openaiApiKey string, state *State) *Controller {
	return &Controller{
		api:          api,
		openAiApiKey: openaiApiKey,
		State:        state,
	}
}

// FailedMsg represents a failure operation.
type FailedMsg struct {
	Err error
}

type UserLoadedMsg struct{ List []list.Item }

// LoadUsers loads users from the AWS API.
func (c *Controller) LoadUsers() tea.Cmd {
	return func() tea.Msg {
		var items []list.Item
		output, err := c.api.ListUsers(context.Background())
		if err != nil {
			return FailedMsg{Err: err}
		}

		for _, user := range output.Users {
			items = append(items, models.User{
				Name: *user.UserName,
				Arn:  *user.Arn,
			})
		}

		return UserLoadedMsg{List: items}
	}
}

type OperationLoadedMsg struct{ List []list.Item }

// LoadOperations loads operations.
func (c *Controller) LoadOperations() tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{
			models.Operation{Id: AttachPolicySlug.String(), Name: ReachableOperations[AttachPolicySlug].Name, Desc: ReachableOperations[AttachPolicySlug].Desc},
			models.Operation{Id: DetachPolicySlug.String(), Name: ReachableOperations[DetachPolicySlug].Name, Desc: ReachableOperations[DetachPolicySlug].Desc},
			models.Operation{Id: AddToGroupSlug.String(), Name: ReachableOperations[AddToGroupSlug].Name, Desc: ReachableOperations[AddToGroupSlug].Desc},
			models.Operation{Id: RemoveFromGroupSlug.String(), Name: ReachableOperations[RemoveFromGroupSlug].Name, Desc: ReachableOperations[RemoveFromGroupSlug].Desc},
			models.Operation{Id: AttachCustomPolicySlug.String(), Name: ReachableOperations[AttachCustomPolicySlug].Name, Desc: ReachableOperations[AttachCustomPolicySlug].Desc},
		}
		return OperationLoadedMsg{List: items}
	}
}

type GroupLoadedMsg struct{ List []list.Item }

// LoadGroups loads groups from the AWS API.
func (c *Controller) LoadGroups() tea.Cmd {
	return func() tea.Msg {
		var items []list.Item

		groups, err := c.api.ListGroups(context.Background())
		if err != nil {
			return FailedMsg{Err: err}
		}

		userGroups, err := c.api.ListGroupsForUser(context.Background(), c.State.GetUser().Name)
		if err != nil {
			return FailedMsg{Err: err}
		}

		switch c.State.operation.Id {
		case AddToGroupSlug.String():
			for _, group := range groups.Groups {
				if !slices.Contains(userGroups, *group.GroupName) {
					items = append(items, models.Group{
						Name: *group.GroupName,
						Arn:  *group.Arn,
					})
				}
			}
		case RemoveFromGroupSlug.String():
			for _, group := range groups.Groups {
				if slices.Contains(userGroups, *group.GroupName) {
					items = append(items, models.Group{
						Name: *group.GroupName,
						Arn:  *group.Arn,
					})
				}
			}
		}

		return GroupLoadedMsg{List: items}
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
		items := []list.Item{
			models.Resource{Name: "All Resources", Arn: "*"},
		}

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

		attachedPolicies, err := c.api.ListAttachedUserPolicies(context.Background(), c.State.GetUser().Name)
		if err != nil {
			return FailedMsg{Err: err}
		}

		switch c.State.operation.Id {
		case AttachPolicySlug.String():
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
		case DetachPolicySlug.String():
			inlinePolicies, err := c.api.ListUserInlinePolicies(context.Background(), c.State.GetUser().Name)
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
			models.PolicyOption{Id: WithoutResourceSlug.String(), Name: ReachablePolicyOptions[WithoutResourceSlug].Name, Desc: ReachablePolicyOptions[WithoutResourceSlug].Desc},
			models.PolicyOption{Id: WithResourceSlug.String(), Name: ReachablePolicyOptions[WithResourceSlug].Name, Desc: ReachablePolicyOptions[WithResourceSlug].Desc},
		}
		return PolicyOptionLoadedMsg{List: items}
	}
}

type OperationType string

// Constants representing user operations and their slugs
const (
	AttachPolicySlug       OperationType = "attach_policy"
	DetachPolicySlug       OperationType = "detach_policy"
	AddToGroupSlug         OperationType = "add_to_group"
	RemoveFromGroupSlug    OperationType = "remove_from_group"
	AttachCustomPolicySlug OperationType = "attach_custom_policy"
)

func (o OperationType) String() string {
	return string(o)
}

// ReachableOperations Predefined list of actions with their names and descriptions
var ReachableOperations = map[OperationType]models.Operation{
	AttachPolicySlug: {
		Id:   AttachPolicySlug.String(),
		Name: "Attach Policy (attach_policy)",
		Desc: "Assign a policy to the user.",
	},
	DetachPolicySlug: {
		Id:   DetachPolicySlug.String(),
		Name: "Detach Policy (detach_policy)",
		Desc: "Remove a policy from the user.",
	},
	AddToGroupSlug: {
		Id:   AddToGroupSlug.String(),
		Name: "Add to Group (add_to_group)",
		Desc: "Include the user in a group.",
	},
	RemoveFromGroupSlug: {
		Id:   RemoveFromGroupSlug.String(),
		Name: "Remove from Group (remove_from_group)",
		Desc: "Exclude the user from a group.",
	},
	AttachCustomPolicySlug: {
		Id:   AttachCustomPolicySlug.String(),
		Name: "Attach Custom Policy (attach_custom_policy)",
		Desc: "Create and attach a custom policy.",
	},
}

type PolicyOptionType string

// Constants representing custom policy options and their slugs
const (
	WithoutResourceSlug PolicyOptionType = "without_resource"
	WithResourceSlug    PolicyOptionType = "with_resource"
)

func (o PolicyOptionType) String() string {
	return string(o)
}

// ReachableCustomPolicyOptions Predefined list of custom policy options with their names and descriptions
var ReachablePolicyOptions = map[PolicyOptionType]models.PolicyOption{
	WithoutResourceSlug: {
		Id:   WithoutResourceSlug.String(),
		Name: "Without Resource (without_resource)",
		Desc: "Applies globally without a resource.",
	},
	WithResourceSlug: {
		Id:   WithResourceSlug.String(),
		Name: "With Resource (with_resource)",
		Desc: "Scoped to a specific resource.",
	},
}

// Next determines the next step based on the current state.
func (c *Controller) Next() tea.Model {
	// Handle case where user is not defined
	if c.State.user == nil {
		return NewUserList(c)
	}

	// Handle case where operation is not defined
	if c.State.operation == nil {
		return NewOperationList(c)
	}

	if c.State.operation.Id == AddToGroupSlug.String() || c.State.operation.Id == RemoveFromGroupSlug.String() {
		// Handle case where group is not defined
		if c.State.GetGroup() == nil {
			return NewGroupList(c)
		}
		return NewResult(c)
	}

	// Handle specific action: AttachCustomPolicySlug
	if c.State.operation.Id == AttachCustomPolicySlug.String() {

		if c.State.policy != nil {
			return NewResult(c)
		}

		// Handle case where a policy option is selected
		if c.State.policyOption != nil {
			switch c.State.policyOption.Id {
			case WithoutResourceSlug.String():
				return NewCreatePolicy(c)

			case WithResourceSlug.String():
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
			// Handle case where resource is defined
			if c.State.resource != nil {
				return NewCreatePolicy(c)
			}

			// Handle case where service is defined
			if c.State.service != nil {
				return NewResourceList(c)
			}
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
	case AttachPolicySlug.String():
		return c.api.AttachPolicyToUser(context.Background(), c.State.GetPolicy().Arn, c.State.GetUser().Name)
	case DetachPolicySlug.String():
		return c.api.DetachPolicyFromUser(context.Background(), c.State.GetPolicy().Arn, c.State.GetUser().Name)
	case AddToGroupSlug.String():
		return c.api.AddUserToGroup(context.Background(), c.State.GetUser().Name, c.State.GetGroup().Name)
	case RemoveFromGroupSlug.String():
		return c.api.RemoveUserFromGroup(context.Background(), c.State.GetUser().Name, c.State.GetGroup().Name)
	case AttachCustomPolicySlug.String():
		output, err := c.api.CreatePolicy(context.Background(), c.State.GetPolicy().Name, c.State.GetPolicy().Document)
		if err != nil {
			return err
		}
		return c.api.AttachPolicyToUser(context.Background(), *output.Policy.Arn, c.State.GetUser().Name)
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
