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

variable "primary_region" {
  description = "The primary AWS region for the global endpoint."
  type        = string
}

variable "secondary_region" {
  description = "The secondary AWS region for failover."
  type        = string
}

variable "event_bus_name" {
  description = "Name of the event bus (must be identical in both regions)."
  type        = string
}

variable "description" {
  description = "Description of the global endpoint."
  type        = string
  default     = null
}

variable "replication_config" {
  description = "Replication configuration. ENABLED requires role_arn."
  type = object({
    state = string
  })
  default = null
}

variable "resource_names_map" {
  description = "Map of resource names for the resource_name module."
  type = map(object({
    name       = string
    max_length = optional(number, 60)
  }))
}

variable "logical_product_family" {
  description = "Logical product family for resource naming."
  type        = string
}

variable "logical_product_service" {
  description = "Logical product service for resource naming."
  type        = string
}

variable "class_env" {
  description = "Class environment for resource naming."
  type        = string
}

variable "instance_env" {
  description = "Instance environment for resource naming."
  type        = number
}

variable "instance_resource" {
  description = "Instance resource for resource naming."
  type        = number
}

variable "tags" {
  description = "Map of tags to assign to resources."
  type        = map(string)
  default     = {}
}
