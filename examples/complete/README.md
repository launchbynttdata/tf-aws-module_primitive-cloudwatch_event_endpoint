# EventBridge Global Endpoint - Complete Example

This example creates an EventBridge global endpoint with:

- Event buses in primary (us-east-1) and secondary (us-west-2) regions
- Route 53 health check for failover triggering
- IAM role for event replication
- Resource naming via the resource_name module

## Usage

```hcl
module "event_endpoint" {
  source = "../.."

  name        = module.resource_names["event_endpoint"].standard
  description = var.description

  event_bus = [
    { event_bus_arn = aws_cloudwatch_event_bus.primary.arn },
    { event_bus_arn = aws_cloudwatch_event_bus.secondary.arn }
  ]

  routing_config = {
    primary_health_check_arn = aws_route53_health_check.endpoint.arn
    secondary_route          = local.secondary_region
  }

  role_arn           = var.replication_config != null && var.replication_config.state == "ENABLED" ? aws_iam_role.replication[0].arn : null
  replication_config = var.replication_config

  tags = var.tags
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| primary_region | The primary AWS region for the global endpoint. | `string` | n/a | yes |
| secondary_region | The secondary AWS region for failover. | `string` | n/a | yes |
| event_bus_name | Name of the event bus (must be identical in both regions). | `string` | n/a | yes |
| description | Description of the global endpoint. | `string` | `null` | no |
| replication_config | Replication configuration. ENABLED requires role_arn. | `object({ state = string })` | `null` | no |
| resource_names_map | Map of resource names for the resource_name module. | `map(object({ name = string, max_length = optional(number, 60) }))` | n/a | yes |
| logical_product_family | Logical product family for resource naming. | `string` | n/a | yes |
| logical_product_service | Logical product service for resource naming. | `string` | n/a | yes |
| class_env | Class environment for resource naming. | `string` | n/a | yes |
| instance_env | Instance environment for resource naming. | `number` | n/a | yes |
| instance_resource | Instance resource for resource naming. | `number` | n/a | yes |
| tags | Map of tags to assign to resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| primary_region | The primary region where the endpoint is created (for API calls). |
| id | The ID of the global endpoint. |
| arn | The ARN of the global endpoint. |
| name | The name of the global endpoint. |
| endpoint_url | The URL of the endpoint used for publishing events. |
| description | The description of the global endpoint. |
| event_bus | The event buses associated with the endpoint. |
| replication_config | The replication configuration of the endpoint. |
| role_arn | The ARN of the IAM role used for replication. |
| routing_config | The routing configuration including failover settings. |

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws.primary"></a> [aws.primary](#provider\_aws.primary) | 5.100.0 |
| <a name="provider_aws.secondary"></a> [aws.secondary](#provider\_aws.secondary) | 5.100.0 |
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_resource_names"></a> [resource\_names](#module\_resource\_names) | terraform.registry.launch.nttdata.com/module_library/resource_name/launch | ~> 2.0 |
| <a name="module_event_endpoint"></a> [event\_endpoint](#module\_event\_endpoint) | ../.. | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_event_bus.primary](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_bus) | resource |
| [aws_cloudwatch_event_bus.secondary](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_bus) | resource |
| [aws_iam_role.replication](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.replication](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_route53_health_check.endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_health_check) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_primary_region"></a> [primary\_region](#input\_primary\_region) | The primary AWS region for the global endpoint. | `string` | n/a | yes |
| <a name="input_secondary_region"></a> [secondary\_region](#input\_secondary\_region) | The secondary AWS region for failover. | `string` | n/a | yes |
| <a name="input_event_bus_name"></a> [event\_bus\_name](#input\_event\_bus\_name) | Name of the event bus (must be identical in both regions). | `string` | n/a | yes |
| <a name="input_description"></a> [description](#input\_description) | Description of the global endpoint. | `string` | `null` | no |
| <a name="input_replication_config"></a> [replication\_config](#input\_replication\_config) | Replication configuration. ENABLED requires role\_arn. | <pre>object({<br/>    state = string<br/>  })</pre> | `null` | no |
| <a name="input_resource_names_map"></a> [resource\_names\_map](#input\_resource\_names\_map) | Map of resource names for the resource\_name module. | <pre>map(object({<br/>    name       = string<br/>    max_length = optional(number, 60)<br/>  }))</pre> | n/a | yes |
| <a name="input_logical_product_family"></a> [logical\_product\_family](#input\_logical\_product\_family) | Logical product family for resource naming. | `string` | n/a | yes |
| <a name="input_logical_product_service"></a> [logical\_product\_service](#input\_logical\_product\_service) | Logical product service for resource naming. | `string` | n/a | yes |
| <a name="input_class_env"></a> [class\_env](#input\_class\_env) | Class environment for resource naming. | `string` | n/a | yes |
| <a name="input_instance_env"></a> [instance\_env](#input\_instance\_env) | Instance environment for resource naming. | `number` | n/a | yes |
| <a name="input_instance_resource"></a> [instance\_resource](#input\_instance\_resource) | Instance resource for resource naming. | `number` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | Map of tags to assign to resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_primary_region"></a> [primary\_region](#output\_primary\_region) | The primary region where the endpoint is created (extracted from ARN for API calls). |
| <a name="output_id"></a> [id](#output\_id) | The ID of the global endpoint. |
| <a name="output_arn"></a> [arn](#output\_arn) | The ARN of the global endpoint. |
| <a name="output_name"></a> [name](#output\_name) | The name of the global endpoint. |
| <a name="output_endpoint_url"></a> [endpoint\_url](#output\_endpoint\_url) | The URL of the endpoint used for publishing events. |
| <a name="output_description"></a> [description](#output\_description) | The description of the global endpoint. |
| <a name="output_event_bus"></a> [event\_bus](#output\_event\_bus) | The event buses associated with the endpoint. |
| <a name="output_replication_config"></a> [replication\_config](#output\_replication\_config) | The replication configuration of the endpoint. |
| <a name="output_role_arn"></a> [role\_arn](#output\_role\_arn) | The ARN of the IAM role used for replication. |
| <a name="output_routing_config"></a> [routing\_config](#output\_routing\_config) | The routing configuration including failover settings. |
<!-- END_TF_DOCS -->
