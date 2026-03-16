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

variable "name" {
  description = "The name of the global endpoint. Maximum of 64 characters consisting of numbers, lower/upper case letters, ., -, _."
  type        = string

  validation {
    condition     = can(regex("^[0-9A-Za-z_.-]{1,64}$", var.name))
    error_message = "Name must be 1-64 characters consisting of numbers, lower/upper case letters, ., -, _."
  }
}

variable "description" {
  description = "A description of the global endpoint."
  type        = string
  default     = null

  validation {
    condition     = var.description == null || (length(var.description) >= 0 && length(var.description) <= 512)
    error_message = "Description must be between 0 and 512 characters."
  }
}

variable "event_bus" {
  description = <<-EOT
    List of exactly two event bus configurations. Each must have event_bus_arn.
    Event bus names must be identical across regions for custom buses.
  EOT
  type = list(object({
    event_bus_arn = string
  }))

  validation {
    condition     = length(var.event_bus) == 2
    error_message = "Exactly two event buses are required (primary and secondary regions)."
  }
}

variable "routing_config" {
  description = <<-EOT
    Routing configuration for failover.
    - primary_health_check_arn: ARN of the Route 53 health check that triggers failover when unhealthy.
    - secondary_route: The secondary region name (e.g., us-west-2) for failover routing.
  EOT
  type = object({
    primary_health_check_arn = string
    secondary_route          = string
  })
}

variable "role_arn" {
  description = "The ARN of the IAM role used for replication between event buses."
  type        = string
  default     = null
}

variable "replication_config" {
  description = <<-EOT
    Replication configuration. When null, replication is not explicitly set (provider default: ENABLED).
    - state: ENABLED or DISABLED. ENABLED replicates events to both regions.
  EOT
  type = object({
    state = string
  })
  default = null

  validation {
    condition     = var.replication_config == null || contains(["ENABLED", "DISABLED"], var.replication_config.state)
    error_message = "Replication state must be ENABLED or DISABLED."
  }
}
