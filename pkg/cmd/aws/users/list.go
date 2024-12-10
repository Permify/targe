package users

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
)

func ListResources(ctx context.Context, client cloudcontrol.ListResourcesAPIClient, resourceTypeName string) ([]string, error) {
	var identifiers []string

	paginator := cloudcontrol.NewListResourcesPaginator(client, &cloudcontrol.ListResourcesInput{
		TypeName: aws.String(resourceTypeName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println(page.ResultMetadata)
		for _, desc := range page.ResourceDescriptions {
			identifiers = append(identifiers, aws.ToString(desc.Identifier))
		}
	}

	return identifiers, nil
}
