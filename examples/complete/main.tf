// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

# Provider aliases for multi-region (EventBridge global endpoints require primary + secondary regions)
# Default provider (for Route53, IAM, endpoint) is provided by provider.tf from Makefile
provider "aws" {
  alias  = "primary"
  region = var.primary_region
}

provider "aws" {
  alias  = "secondary"
  region = var.secondary_region
}

data "aws_caller_identity" "current" {}

locals {
  primary_region   = var.primary_region
  secondary_region = var.secondary_region
  account_id       = data.aws_caller_identity.current.account_id

  primary_event_bus_arn   = "arn:aws:events:${local.primary_region}:${local.account_id}:event-bus/${var.event_bus_name}"
  secondary_event_bus_arn = "arn:aws:events:${local.secondary_region}:${local.account_id}:event-bus/${var.event_bus_name}"

  replication_enabled = var.replication_config != null ? var.replication_config.state == "ENABLED" : false
}

# Resource naming
module "resource_names" {
  source   = "terraform.registry.launch.nttdata.com/module_library/resource_name/launch"
  version  = "~> 2.0"
  for_each = var.resource_names_map

  logical_product_family  = var.logical_product_family
  logical_product_service = var.logical_product_service
  class_env               = var.class_env
  instance_env            = var.instance_env
  instance_resource       = var.instance_resource
  cloud_resource_type     = each.value.name
  maximum_length          = each.value.max_length

  # AWS regions have hyphens (e.g. "us-east-1") — strip them for resource naming
  region = join("", split("-", local.primary_region))
}

# Event buses in primary and secondary regions (identical names required)
resource "aws_cloudwatch_event_bus" "primary" {
  provider = aws.primary

  name = var.event_bus_name
  tags = var.tags
}

resource "aws_cloudwatch_event_bus" "secondary" {
  provider = aws.secondary

  name = var.event_bus_name
  tags = var.tags
}

# Route 53 health check for failover triggering
# Uses a reliable HTTPS endpoint; for production, use a CloudWatch metric health check
# that monitors EventBridge IngestionToInvocationStartLatency
resource "aws_route53_health_check" "endpoint" {
  fqdn              = "aws.amazon.com"
  port              = 443
  type              = "HTTPS"
  resource_path     = "/"
  request_interval  = 30
  failure_threshold = 3

  measure_latency = true

  tags = var.tags
}

# IAM role for event replication (required when replication is ENABLED)
resource "aws_iam_role" "replication" {
  count = local.replication_enabled ? 1 : 0

  name = module.resource_names["iam_role"].standard

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "events.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "replication" {
  count = local.replication_enabled ? 1 : 0

  name = "${module.resource_names["iam_role"].standard}-replication"
  role = aws_iam_role.replication[0].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "events:PutEvents"
        Resource = [
          local.primary_event_bus_arn,
          local.secondary_event_bus_arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "events:PutRule",
          "events:PutTargets",
          "events:DeleteRule",
          "events:RemoveTargets",
          "events:DescribeRule",
          "events:DescribeEventBus"
        ]
        Resource = [
          local.primary_event_bus_arn,
          local.secondary_event_bus_arn,
          "arn:aws:events:${local.primary_region}:${local.account_id}:rule/${var.event_bus_name}/*",
          "arn:aws:events:${local.secondary_region}:${local.account_id}:rule/${var.event_bus_name}/*"
        ]
      },
      {
        Effect   = "Allow"
        Action   = "iam:PassRole"
        Resource = aws_iam_role.replication[0].arn
        Condition = {
          StringEquals = {
            "iam:PassedToService" = "events.amazonaws.com"
          }
        }
      }
    ]
  })
}

# EventBridge Global Endpoint (must be created in primary region)
module "event_endpoint" {
  source = "../.."

  providers = {
    aws = aws.primary
  }

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

  role_arn           = local.replication_enabled ? aws_iam_role.replication[0].arn : null
  replication_config = var.replication_config
}
