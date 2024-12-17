package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func ListResourcest(ctx context.Context, client cloudcontrol.ListResourcesAPIClient, resourceTypeName string) ([]string, error) {
	var identifiers []string

	paginator := cloudcontrol.NewListResourcesPaginator(client, &cloudcontrol.ListResourcesInput{
		TypeName: aws.String(resourceTypeName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, desc := range page.ResourceDescriptions {
			identifiers = append(identifiers, aws.ToString(desc.Identifier))
			if desc.Properties != nil {
				fmt.Println(*desc.Properties)
			}
		}
	}

	return identifiers, nil
}

func ListAttachedUserPolicies(ctx context.Context, client *iam.Client, userName string) ([]string, error) {
	paginator := iam.NewListAttachedUserPoliciesPaginator(client, &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(userName),
	})

	var allPolicyNames []string

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list attached user policies: %v", err)
		}

		for _, policy := range output.AttachedPolicies {
			allPolicyNames = append(allPolicyNames, aws.ToString(policy.PolicyName))
		}
	}

	return allPolicyNames, nil
}

func ListUserPolicies(ctx context.Context, client *iam.Client, userName string) ([]string, error) {
	paginator := iam.NewListUserPoliciesPaginator(client, &iam.ListUserPoliciesInput{
		UserName: aws.String(userName),
	})

	var allPolicies []string

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list user policies: %v", err)
		}

		allPolicies = append(allPolicies, output.PolicyNames...)
	}

	return allPolicies, nil
}

func ListGroupsForUser(ctx context.Context, client *iam.Client, userName string) ([]string, error) {
	// Create the paginator for ListGroupsForUser
	paginator := iam.NewListGroupsForUserPaginator(client, &iam.ListGroupsForUserInput{
		UserName: aws.String(userName),
	})

	var groupNames []string

	// Iterate through all pages
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list groups for user: %v", err)
		}

		// Extract group names and append to the slice
		for _, group := range output.Groups {
			groupNames = append(groupNames, *group.GroupName)
		}
	}

	return groupNames, nil
}

func ListGroupPolicies(ctx context.Context, client *iam.Client, groupName string) ([]string, error) {
	// Create the paginator for ListGroupPolicies
	paginator := iam.NewListGroupPoliciesPaginator(client, &iam.ListGroupPoliciesInput{
		GroupName: aws.String(groupName),
	})

	var inlinePolicies []string

	// Iterate through all pages
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list inline group policies: %v", err)
		}

		// Append policy names
		inlinePolicies = append(inlinePolicies, output.PolicyNames...)
	}

	return inlinePolicies, nil
}

func ListAttachedGroupPolicies(ctx context.Context, client *iam.Client, groupName string) ([]string, error) {
	// Create the paginator for ListAttachedGroupPolicies
	paginator := iam.NewListAttachedGroupPoliciesPaginator(client, &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(groupName),
	})

	var attachedPolicies []string

	// Iterate through all pages
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list attached group policies: %v", err)
		}

		// Append attached policy names
		for _, policy := range output.AttachedPolicies {
			attachedPolicies = append(attachedPolicies, *policy.PolicyName)
		}
	}

	return attachedPolicies, nil
}

func ListAllUserGroupPolicies(ctx context.Context, client *iam.Client, userName string) ([]string, error) {
	// Fetch all groups for the user
	groupNames, err := ListGroupsForUser(ctx, client, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups for user %s: %v", userName, err)
	}

	allUserGroupPolicies := []string{}
	// Iterate through all groups
	for _, groupName := range groupNames {
		inlinePolicies, err := ListGroupPolicies(ctx, client, groupName)
		if err != nil {
			return nil, fmt.Errorf("failed to list inline policies for group %s: %v", groupName, err)
		}

		allUserGroupPolicies = append(allUserGroupPolicies, inlinePolicies...)

		attachedPolicies, err := ListAttachedGroupPolicies(ctx, client, groupName)
		if err != nil {
			return nil, fmt.Errorf("failed to list attached policies for group %s: %v", groupName, err)
		}

		allUserGroupPolicies = append(allUserGroupPolicies, attachedPolicies...)
	}

	return allUserGroupPolicies, nil
}

func ListAllUserPolicies(ctx context.Context, client *iam.Client, userName string) ([]string, error) {
	inlineUserPolicies, err := ListUserPolicies(ctx, client, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to list user policies: %v", err)
	}
	allPolicies := append([]string{}, inlineUserPolicies...)

	attachedUserPolicies, err := ListAttachedUserPolicies(ctx, client, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to list attached user policies: %v", err)
	}
	allPolicies = append(allPolicies, attachedUserPolicies...)

	userGroupPolicies, err := ListAllUserGroupPolicies(ctx, client, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to list user group policies: %v", err)
	}

	allPolicies = append(allPolicies, userGroupPolicies...)

	return allPolicies, nil
}

// Define the policy structure
type Policy struct {
	Name string `json:"name"`
	Arn  string `json:"arn"`
}

func ListAWSManagedPolicies(ctx context.Context, cfg aws.Config) ([]Policy, error) {
	var allPolicies []Policy

	client := iam.NewFromConfig(cfg)
	paginator := iam.NewListPoliciesPaginator(client, &iam.ListPoliciesInput{
		Scope: types.PolicyScopeType("AWS"),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list policies: %v", err)
		}

		for _, policy := range output.Policies {
			allPolicies = append(allPolicies, Policy{
				Name: *policy.PolicyName,
				Arn:  *policy.Arn,
			})
		}
	}

	return allPolicies, nil
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}

	awsManagedPolicies, _ := ListAWSManagedPolicies(ctx, cfg)
	// Convert to JSON
	jsonData, err := json.MarshalIndent(awsManagedPolicies, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write to file
	err = os.WriteFile("policies.json", jsonData, 0o644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
