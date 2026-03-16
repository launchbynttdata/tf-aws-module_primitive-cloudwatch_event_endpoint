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

resource "aws_cloudwatch_event_endpoint" "endpoint" {
  name        = var.name
  description = var.description

  dynamic "event_bus" {
    for_each = var.event_bus
    content {
      event_bus_arn = event_bus.value.event_bus_arn
    }
  }

  routing_config {
    failover_config {
      primary {
        health_check = var.routing_config.primary_health_check_arn
      }
      secondary {
        route = var.routing_config.secondary_route
      }
    }
  }

  role_arn = var.role_arn

  dynamic "replication_config" {
    for_each = var.replication_config != null ? [var.replication_config] : []
    content {
      state = replication_config.value.state
    }
  }
}
