# tf-aws-module_primitive-cloudwatch_event_endpoint

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

This Terraform module creates an [AWS EventBridge Global Endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_endpoint) for regional fault tolerance. Global endpoints enable automatic failover between primary and secondary regions when the Route 53 health check reports the primary region as unhealthy.

## Features

- Creates an EventBridge global endpoint with configurable routing
- Supports event replication between primary and secondary regions
- Configurable Route 53 health check for failover triggering
- IAM role for event replication (when replication is enabled)
- Full support for all documented resource attributes

## Usage

```hcl
module "event_endpoint" {
  source = "terraform.registry.launch.nttdata.com/module_primitive/cloudwatch_event_endpoint/aws"
  version = "~> 1.0"

  name        = "my-global-endpoint"
  description = "EventBridge global endpoint for regional fault tolerance"

  event_bus = [
    { event_bus_arn = aws_cloudwatch_event_bus.primary.arn },
    { event_bus_arn = aws_cloudwatch_event_bus.secondary.arn }
  ]

  routing_config = {
    primary_health_check_arn = aws_route53_health_check.endpoint.arn
    secondary_route          = "us-west-2"
  }

  role_arn = aws_iam_role.replication.arn

  replication_config = {
    state = "ENABLED"
  }
}
```

## Pre-Commit Hooks

The [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines hooks for Terraform, Go, and common linting. The `commitlint` hook enforces conventional commit messages.

The `detect-secrets-hook` prevents new secrets from being introduced. See [detect-secrets documentation](https://github.com/Yelp/detect-secrets) for details.

To install the commit-msg hook:

```shell
pre-commit install --hook-type commit-msg
```

## Local Testing

1. Install components: `make configure`
2. Set up AWS credentials (e.g., `AWS_PROFILE` or `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`)
3. Create `examples/complete/provider.tf` with AWS provider configuration (or rely on Makefile-generated file)
4. Create `examples/complete/terraform.tfvars` with variable values
5. Run validation: `make check`

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_event_endpoint.endpoint](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_event_endpoint) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_description"></a> [description](#input\_description) | A description of the global endpoint. | `string` | `null` | no |
| <a name="input_event_bus"></a> [event\_bus](#input\_event\_bus) | List of exactly two event bus configurations. Each must have event\_bus\_arn.<br/>Event bus names must be identical across regions for custom buses. | <pre>list(object({<br/>    event_bus_arn = string<br/>  }))</pre> | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | The name of the global endpoint. Maximum of 64 characters consisting of numbers, lower/upper case letters, ., -, \_. | `string` | n/a | yes |
| <a name="input_replication_config"></a> [replication\_config](#input\_replication\_config) | Replication configuration. When null, replication is not explicitly set (provider default: ENABLED).<br/>- state: ENABLED or DISABLED. ENABLED replicates events to both regions. | <pre>object({<br/>    state = string<br/>  })</pre> | `null` | no |
| <a name="input_role_arn"></a> [role\_arn](#input\_role\_arn) | The ARN of the IAM role used for replication between event buses. | `string` | `null` | no |
| <a name="input_routing_config"></a> [routing\_config](#input\_routing\_config) | Routing configuration for failover.<br/>- primary\_health\_check\_arn: ARN of the Route 53 health check that triggers failover when unhealthy.<br/>- secondary\_route: The secondary region name (e.g., us-west-2) for failover routing. | <pre>object({<br/>    primary_health_check_arn = string<br/>    secondary_route          = string<br/>  })</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_arn"></a> [arn](#output\_arn) | The ARN of the global endpoint. |
| <a name="output_description"></a> [description](#output\_description) | The description of the global endpoint. |
| <a name="output_endpoint_url"></a> [endpoint\_url](#output\_endpoint\_url) | The URL of the endpoint used for publishing events. |
| <a name="output_event_bus"></a> [event\_bus](#output\_event\_bus) | The event buses associated with the endpoint. |
| <a name="output_id"></a> [id](#output\_id) | The ID of the resource (same as the name). |
| <a name="output_name"></a> [name](#output\_name) | The name of the global endpoint. |
| <a name="output_replication_config"></a> [replication\_config](#output\_replication\_config) | The replication configuration of the endpoint. |
| <a name="output_role_arn"></a> [role\_arn](#output\_role\_arn) | The ARN of the IAM role used for replication. |
| <a name="output_routing_config"></a> [routing\_config](#output\_routing\_config) | The routing configuration including failover settings. |
<!-- END_TF_DOCS -->
