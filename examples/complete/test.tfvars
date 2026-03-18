primary_region   = "us-east-1"
secondary_region = "us-west-2"
event_bus_name   = "example-global-endpoint-bus"

description = "Example EventBridge global endpoint for regional fault tolerance"

replication_config = {
  state = "ENABLED"
}

resource_names_map = {
  event_endpoint = {
    name       = "ep"
    max_length = 64
  }
  iam_role = {
    name       = "role"
    max_length = 64
  }
}

logical_product_family  = "launch"
logical_product_service = "events"
class_env               = "dev"
instance_env            = 0
instance_resource       = 1

tags = {
  Environment = "example"
  ManagedBy   = "terraform"
}
