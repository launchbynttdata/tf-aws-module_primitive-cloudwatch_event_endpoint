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

output "id" {
  description = "The ID of the resource (same as the name)."
  value       = aws_cloudwatch_event_endpoint.endpoint.id
}

output "arn" {
  description = "The ARN of the global endpoint."
  value       = aws_cloudwatch_event_endpoint.endpoint.arn
}

output "name" {
  description = "The name of the global endpoint."
  value       = aws_cloudwatch_event_endpoint.endpoint.name
}

output "endpoint_url" {
  description = "The URL of the endpoint used for publishing events."
  value       = aws_cloudwatch_event_endpoint.endpoint.endpoint_url
}

output "description" {
  description = "The description of the global endpoint."
  value       = aws_cloudwatch_event_endpoint.endpoint.description
}

output "event_bus" {
  description = "The event buses associated with the endpoint."
  value       = aws_cloudwatch_event_endpoint.endpoint.event_bus
}

output "replication_config" {
  description = "The replication configuration of the endpoint."
  value       = aws_cloudwatch_event_endpoint.endpoint.replication_config
}

output "role_arn" {
  description = "The ARN of the IAM role used for replication."
  value       = aws_cloudwatch_event_endpoint.endpoint.role_arn
}

output "routing_config" {
  description = "The routing configuration including failover settings."
  value       = aws_cloudwatch_event_endpoint.endpoint.routing_config
}
