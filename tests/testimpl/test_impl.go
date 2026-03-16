package testimpl

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEventBridgeClient(t *testing.T, opts *terraform.Options) *eventbridge.Client {
	region := terraform.Output(t, opts, "primary_region")
	require.NotEmpty(t, region, "primary_region output required for API calls")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	require.NoError(t, err, "unable to load AWS config for region %s", region)
	return eventbridge.NewFromConfig(cfg)
}

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		id := terraform.Output(t, opts, "id")
		name := terraform.Output(t, opts, "name")
		arn := terraform.Output(t, opts, "arn")
		endpointURL := terraform.Output(t, opts, "endpoint_url")

		assert.Equal(t, name, id, "id should equal name for EventBridge global endpoint")
		require.NotEmpty(t, arn, "arn should be set")
		require.NotEmpty(t, endpointURL, "endpoint_url should be set")
	})

	t.Run("VerifyEndpointViaAWSAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.Output(t, opts, "name")
		expectedARN := terraform.Output(t, opts, "arn")

		client := getEventBridgeClient(t, opts)
		output, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint should succeed")
		require.NotNil(t, output, "DescribeEndpoint output should not be nil")

		assert.Equal(t, endpointName, aws.ToString(output.Name), "endpoint name should match")
		assert.Equal(t, expectedARN, aws.ToString(output.Arn), "endpoint ARN should match Terraform output")
		assert.Equal(t, "ACTIVE", string(output.State), "endpoint state should be ACTIVE")
		require.Len(t, output.EventBuses, 2, "endpoint should have exactly two event buses")
	})

	t.Run("VerifyReplicationConfigViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.Output(t, opts, "name")
		expectedRoleARN := terraform.Output(t, opts, "role_arn")

		client := getEventBridgeClient(t, opts)
		output, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint should succeed")

		if expectedRoleARN != "" {
			require.NotNil(t, output.ReplicationConfig, "replication config should be present when role is set")
			assert.Equal(t, "ENABLED", string(output.ReplicationConfig.State), "replication state should be ENABLED")
			require.NotNil(t, output.RoleArn, "role ARN should be set when replication is enabled")
			assert.Equal(t, expectedRoleARN, aws.ToString(output.RoleArn), "role ARN should match Terraform output")
		}
	})

	t.Run("PutEventsToEndpoint", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.Output(t, opts, "name")
		require.NotEmpty(t, endpointName, "endpoint name must be set for PutEvents test")

		client := getEventBridgeClient(t, opts)
		desc, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint must succeed to obtain EndpointId for PutEvents")
		endpointID := aws.ToString(desc.EndpointId)
		require.NotEmpty(t, endpointID, "DescribeEndpoint must return EndpointId for PutEvents")

		// Send event to the global endpoint (write operation)
		// Event bus name from example test.tfvars: example-global-endpoint-bus
		eventBusName := "example-global-endpoint-bus"

		result, err := client.PutEvents(context.Background(), &eventbridge.PutEventsInput{
			EndpointId: aws.String(endpointID),
			Entries: []eventbridgetypes.PutEventsRequestEntry{
				{
					Source:       aws.String("test.terraform"),
					DetailType:   aws.String("Terraform Functional Test"),
					Detail:       aws.String(`{"test": "functional","module": "cloudwatch_event_endpoint"}`),
					EventBusName: aws.String(eventBusName),
				},
			},
		})
		require.NoError(t, err, "PutEvents to global endpoint should succeed")
		require.NotEmpty(t, result.Entries, "PutEvents should return entries")
		require.Nil(t, result.Entries[0].ErrorCode, "PutEvents failed: %s", aws.ToString(result.Entries[0].ErrorMessage))
	})
}

func TestComposableCompleteReadonly(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputsReadonly", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		id := terraform.Output(t, opts, "id")
		name := terraform.Output(t, opts, "name")
		arn := terraform.Output(t, opts, "arn")

		assert.Equal(t, name, id, "id should equal name for EventBridge global endpoint")
		require.NotEmpty(t, arn, "arn should be set")
	})

	t.Run("VerifyEndpointExistsViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.Output(t, opts, "name")

		client := getEventBridgeClient(t, opts)
		output, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint should succeed")
		require.NotNil(t, output, "DescribeEndpoint output should not be nil")

		assert.Equal(t, endpointName, aws.ToString(output.Name), "endpoint name should match")
		assert.Equal(t, "ACTIVE", string(output.State), "endpoint state should be ACTIVE")
		require.Len(t, output.EventBuses, 2, "endpoint should have exactly two event buses")
	})

	t.Run("VerifyRoutingConfigViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.Output(t, opts, "name")

		client := getEventBridgeClient(t, opts)
		output, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint should succeed")
		require.NotNil(t, output.RoutingConfig, "routing config should be present")
		require.NotNil(t, output.RoutingConfig.FailoverConfig, "failover config should be present")
		require.NotNil(t, output.RoutingConfig.FailoverConfig.Secondary, "secondary config should be present")

		secondaryRoute := aws.ToString(output.RoutingConfig.FailoverConfig.Secondary.Route)
		require.NotEmpty(t, secondaryRoute, "secondary route should be set")
	})
}
