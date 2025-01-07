package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type Api struct {
	client *iam.Client
	config aws.Config
}

func NewApi(config aws.Config) *Api {
	return &Api{
		client: iam.NewFromConfig(config),
		config: config,
	}
}

func (op *Api) CreatePolicy(ctx context.Context, name, document string) (*iam.CreatePolicyOutput, error) {
	return op.client.CreatePolicy(ctx, &iam.CreatePolicyInput{
		Description:    aws.String("created by targe"),
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(document),
	})
}

func (op *Api) AttachPolicyToUser(ctx context.Context, policyArn, username string) error {
	_, err := op.client.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
		PolicyArn: aws.String(policyArn),
		UserName:  aws.String(username),
	})
	return err
}

func (op *Api) DetachPolicyFromUser(ctx context.Context, policyArn, username string) error {
	_, err := op.client.DetachUserPolicy(ctx, &iam.DetachUserPolicyInput{
		PolicyArn: aws.String(policyArn),
		UserName:  aws.String(username),
	})
	return err
}

func (op *Api) AttachPolicyToGroup(ctx context.Context, policyArn, groupname string) error {
	_, err := op.client.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
		PolicyArn: aws.String(policyArn),
		GroupName: aws.String(groupname),
	})
	return err
}

func (op *Api) DetachPolicyFromGroup(ctx context.Context, policyArn, groupname string) error {
	_, err := op.client.DetachGroupPolicy(ctx, &iam.DetachGroupPolicyInput{
		PolicyArn: aws.String(policyArn),
		GroupName: aws.String(groupname),
	})
	return err
}

func (op *Api) AttachPolicyToRole(ctx context.Context, policyArn, rolename string) error {
	_, err := op.client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(rolename),
	})
	return err
}

func (op *Api) DetachPolicyFromRole(ctx context.Context, policyArn, rolename string) error {
	_, err := op.client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(rolename),
	})
	return err
}

func (op *Api) PutInlinePolicyToUser(ctx context.Context, policyname, policyDocument, username string) error {
	_, err := op.client.PutUserPolicy(ctx, &iam.PutUserPolicyInput{
		PolicyName:     aws.String(policyname),
		PolicyDocument: aws.String(policyDocument),
		UserName:       aws.String(username),
	})
	return err
}

func (op *Api) DeleteInlinePolicyFromUser(ctx context.Context, policyname, username string) error {
	_, err := op.client.DeleteUserPolicy(ctx, &iam.DeleteUserPolicyInput{
		PolicyName: aws.String(policyname),
		UserName:   aws.String(username),
	})
	return err
}

func (op *Api) PutInlinePolicyToGroup(ctx context.Context, policyname, policyDocument, groupname string) error {
	_, err := op.client.PutGroupPolicy(ctx, &iam.PutGroupPolicyInput{
		PolicyName:     aws.String(policyname),
		PolicyDocument: aws.String(policyDocument),
		GroupName:      aws.String(groupname),
	})
	return err
}

func (op *Api) DeleteInlinePolicyFromGroup(ctx context.Context, policyname, groupname string) error {
	_, err := op.client.DeleteGroupPolicy(ctx, &iam.DeleteGroupPolicyInput{
		PolicyName: aws.String(policyname),
		GroupName:  aws.String(groupname),
	})
	return err
}

func (op *Api) PutInlinePolicyToRole(ctx context.Context, policyname, policyDocument, rolename string) error {
	_, err := op.client.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		PolicyName:     aws.String(policyname),
		PolicyDocument: aws.String(policyDocument),
		RoleName:       aws.String(rolename),
	})
	return err
}

func (op *Api) DeleteInlinePolicyFromRole(ctx context.Context, policyname, rolename string) error {
	_, err := op.client.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(policyname),
		RoleName:   aws.String(rolename),
	})
	return err
}

func (op *Api) AddUserToGroup(ctx context.Context, username, groupname string) error {
	_, err := op.client.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		GroupName: aws.String(groupname),
		UserName:  aws.String(username),
	})
	return err
}

func (op *Api) RemoveUserFromGroup(ctx context.Context, username, groupname string) error {
	_, err := op.client.RemoveUserFromGroup(ctx, &iam.RemoveUserFromGroupInput{
		GroupName: aws.String(groupname),
		UserName:  aws.String(username),
	})
	return err
}
