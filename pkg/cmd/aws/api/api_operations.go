package api

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func AttachPolicyToUser(ctx context.Context, cfg aws.Config, policyArn, userName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyArn == "" || userName == "" {
		return fmt.Errorf("policyArn and userName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
		PolicyArn: aws.String(policyArn),
		UserName:  aws.String(userName),
	})

	return err
}

func DetachPolicyFromUser(ctx context.Context, cfg aws.Config, policyArn, userName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyArn == "" || userName == "" {
		return fmt.Errorf("policyArn and userName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.DetachUserPolicy(ctx, &iam.DetachUserPolicyInput{
		PolicyArn: aws.String(policyArn),
		UserName:  aws.String(userName),
	})

	return err
}

func AttachPolicyToGroup(ctx context.Context, cfg aws.Config, policyArn, groupName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyArn == "" || groupName == "" {
		return fmt.Errorf("policyArn and groupName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
		PolicyArn: aws.String(policyArn),
		GroupName: aws.String(groupName),
	})

	return err
}

func DetachPolicyFromGroup(ctx context.Context, cfg aws.Config, policyArn, groupName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyArn == "" || groupName == "" {
		return fmt.Errorf("policyArn and groupName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.DetachGroupPolicy(ctx, &iam.DetachGroupPolicyInput{
		PolicyArn: aws.String(policyArn),
		GroupName: aws.String(groupName),
	})

	return err
}

func AttachPolicyToRole(ctx context.Context, cfg aws.Config, policyArn, roleName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyArn == "" || roleName == "" {
		return fmt.Errorf("policyArn and roleName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	return err
}

func DetachPolicyFromRole(ctx context.Context, cfg aws.Config, policyArn, roleName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyArn == "" || roleName == "" {
		return fmt.Errorf("policyArn and roleName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	return err
}

func PutInlinePolicyToUser(ctx context.Context, cfg aws.Config, policyName, policyDocument, userName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyName == "" || policyDocument == "" || userName == "" {
		return fmt.Errorf("policyName, policyDocument and userName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.PutUserPolicy(ctx, &iam.PutUserPolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
		UserName:       aws.String(userName),
	})

	return err
}

func DeleteInlinePolicyFromUser(ctx context.Context, cfg aws.Config, policyName, userName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyName == "" || userName == "" {
		return fmt.Errorf("policyName and userName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.DeleteUserPolicy(ctx, &iam.DeleteUserPolicyInput{
		PolicyName: aws.String(policyName),
		UserName:   aws.String(userName),
	})

	return err
}

func PutInlinePolicyToGroup(ctx context.Context, cfg aws.Config, policyName, policyDocument, groupName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyName == "" || policyDocument == "" || groupName == "" {
		return fmt.Errorf("policyName, policyDocument and groupName cannot be empty")
	}
	client := iam.NewFromConfig(cfg)
	_, err := client.PutGroupPolicy(ctx, &iam.PutGroupPolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
		GroupName:      aws.String(groupName),
	})

	return err
}

func DeleteInlinePolicyFromGroup(ctx context.Context, cfg aws.Config, policyName, groupName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyName == "" || groupName == "" {
		return fmt.Errorf("policyName and groupName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.DeleteGroupPolicy(ctx, &iam.DeleteGroupPolicyInput{
		PolicyName: aws.String(policyName),
		GroupName:  aws.String(groupName),
	})

	return err
}

func PutInlinePolicyToRole(ctx context.Context, cfg aws.Config, policyName, policyDocument, roleName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyName == "" || policyDocument == "" || roleName == "" {
		return fmt.Errorf("policyName, policyDocument and roleName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
		RoleName:       aws.String(roleName),
	})

	return err

}

func DeleteInlinePolicyFromRole(ctx context.Context, cfg aws.Config, policyName, roleName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if policyName == "" || roleName == "" {
		return fmt.Errorf("policyName and roleName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	})

	return err
}

func AddUserToGroup(ctx context.Context, cfg aws.Config, userName, groupName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if userName == "" || groupName == "" {
		return fmt.Errorf("userName and groupName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		GroupName: aws.String(groupName),
		UserName:  aws.String(userName),
	})

	return err

}

func RemoveUserFromGroup(ctx context.Context, cfg aws.Config, userName, groupName string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if userName == "" || groupName == "" {
		return fmt.Errorf("userName and groupName cannot be empty")
	}

	client := iam.NewFromConfig(cfg)
	_, err := client.RemoveUserFromGroup(ctx, &iam.RemoveUserFromGroupInput{
		GroupName: aws.String(groupName),
		UserName:  aws.String(userName),
	})

	return err
}
