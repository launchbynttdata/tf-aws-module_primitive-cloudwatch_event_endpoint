package testimpl

import (
	"context"
	"net/url"
	"regexp"
	"strings"
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

var eventEndpointARNPattern = regexp.MustCompile(`^arn:aws[^:]*:events:[a-z0-9-]+:[0-9]{12}:endpoint/.+$`)

func getEventBridgeClient(t *testing.T, opts *terraform.Options) *eventbridge.Client {
	region := terraform.OutputContext(t, context.Background(), opts, "primary_region")
	require.NotEmpty(t, region, "primary_region output required for API calls")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	require.NoError(t, err, "unable to load AWS config for region %s", region)
	return eventbridge.NewFromConfig(cfg)
}

func eventBusNameFromARN(t *testing.T, eventBusARN string) string {
	t.Helper()
	require.NotEmpty(t, eventBusARN, "event bus ARN output must be set")

	separatorIdx := strings.LastIndex(eventBusARN, "/")
	require.Greater(t, separatorIdx, 0, "event bus ARN must contain resource name after '/'")

	eventBusName := eventBusARN[separatorIdx+1:]
	require.NotEmpty(t, eventBusName, "event bus ARN must include a non-empty event bus name")
	return eventBusName
}

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		id := terraform.OutputContext(t, context.Background(), opts, "id")
		name := terraform.OutputContext(t, context.Background(), opts, "name")
		arn := terraform.OutputContext(t, context.Background(), opts, "arn")
		endpointURL := terraform.OutputContext(t, context.Background(), opts, "endpoint_url")

		assert.Equal(t, name, id, "id should equal name for EventBridge global endpoint")
		require.True(t, eventEndpointARNPattern.MatchString(arn), "arn should match EventBridge endpoint ARN format")

		parsedEndpointURL, err := url.Parse(endpointURL)
		require.NoError(t, err, "endpoint_url should be a valid URL")
		assert.Equal(t, "https", parsedEndpointURL.Scheme, "endpoint_url should use HTTPS")
		require.NotEmpty(t, parsedEndpointURL.Host, "endpoint_url should include a host")
	})

	t.Run("VerifyEndpointViaAWSAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.OutputContext(t, context.Background(), opts, "name")
		expectedARN := terraform.OutputContext(t, context.Background(), opts, "arn")
		expectedEndpointURL := terraform.OutputContext(t, context.Background(), opts, "endpoint_url")

		client := getEventBridgeClient(t, opts)
		output, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint should succeed")
		require.NotNil(t, output, "DescribeEndpoint output should not be nil")

		assert.Equal(t, endpointName, aws.ToString(output.Name), "endpoint name should match")
		assert.Equal(t, expectedARN, aws.ToString(output.Arn), "endpoint ARN should match Terraform output")
		assert.Equal(t, expectedEndpointURL, aws.ToString(output.EndpointUrl), "endpoint URL should match Terraform output")
		assert.Equal(t, "ACTIVE", string(output.State), "endpoint state should be ACTIVE")
		require.Len(t, output.EventBuses, 2, "endpoint should have exactly two event buses")
	})

	t.Run("VerifyReplicationConfigViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.OutputContext(t, context.Background(), opts, "name")
		expectedRoleARN := terraform.OutputContext(t, context.Background(), opts, "role_arn")

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
		endpointName := terraform.OutputContext(t, context.Background(), opts, "name")
		require.NotEmpty(t, endpointName, "endpoint name must be set for PutEvents test")

		client := getEventBridgeClient(t, opts)
		desc, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint must succeed to obtain EndpointId for PutEvents")
		endpointID := aws.ToString(desc.EndpointId)
		require.NotEmpty(t, endpointID, "DescribeEndpoint must return EndpointId for PutEvents")

		eventBuses := terraform.OutputListOfObjectsContext(t, context.Background(), opts, "event_bus")
		require.Len(t, eventBuses, 2, "event_bus output should contain two buses")

		eventBusARNValue, ok := eventBuses[0]["event_bus_arn"]
		require.True(t, ok, "event_bus output should include event_bus_arn")

		eventBusARN, ok := eventBusARNValue.(string)
		require.True(t, ok, "event_bus_arn output should be a string")

		eventBusName := eventBusNameFromARN(t, eventBusARN)

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
		id := terraform.OutputContext(t, context.Background(), opts, "id")
		name := terraform.OutputContext(t, context.Background(), opts, "name")
		arn := terraform.OutputContext(t, context.Background(), opts, "arn")

		assert.Equal(t, name, id, "id should equal name for EventBridge global endpoint")
		require.True(t, eventEndpointARNPattern.MatchString(arn), "arn should match EventBridge endpoint ARN format")
	})

	t.Run("VerifyEndpointExistsViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		endpointName := terraform.OutputContext(t, context.Background(), opts, "name")

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
		endpointName := terraform.OutputContext(t, context.Background(), opts, "name")

		client := getEventBridgeClient(t, opts)
		output, err := client.DescribeEndpoint(context.Background(), &eventbridge.DescribeEndpointInput{
			Name: aws.String(endpointName),
		})
		require.NoError(t, err, "DescribeEndpoint should succeed")
		require.NotNil(t, output.RoutingConfig, "routing config should be present")
		require.NotNil(t, output.RoutingConfig.FailoverConfig, "failover config should be present")
		require.NotNil(t, output.RoutingConfig.FailoverConfig.Secondary, "secondary config should be present")

		secondaryRoute := aws.ToString(output.RoutingConfig.FailoverConfig.Secondary.Route)
		expectedRoutingConfig := terraform.OutputMapContext(t, context.Background(), opts, "routing_config")
		expectedSecondaryRoute := expectedRoutingConfig["secondary_route"]

		require.NotEmpty(t, expectedSecondaryRoute, "routing_config.secondary_route output should be set")
		assert.Equal(t, expectedSecondaryRoute, secondaryRoute, "secondary route should match Terraform output")
	})
}
