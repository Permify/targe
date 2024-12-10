package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
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

func main() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := cloudcontrol.NewFromConfig(cfg)

	resources, err := ListResourcest(ctx, client, "AWS::RDS::DBCluster")
	if err != nil {
		panic(err)
	}

	for _, r := range resources {
		fmt.Println(r)
	}
}
