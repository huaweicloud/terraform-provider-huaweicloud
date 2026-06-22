service_name        = "tf_test_service"
service_description = "Created online inference service by Terraform"
service_version     = "0.0.1"
service_deploy_type = "MULTI"

service_group_name    = "tf_test_deploy_group"
service_group_pool_id = "pool-7b504f98-914a-4a03-93f1-d093f7b27908"

service_unit_configs = [
  {
    count    = 1
    recovery = "INSTANCE"
    role     = "COMMON"
    flavor   = "your_unit_resource_flavor" # Please replace with your own actual resource flavor.

    image = {
      source   = "SWR"
      swr_path = "your_swr_image_path" # Please replace with your own actual SWR path.
    }
    readiness_health = {
      initial_delay_seconds = 60
      timeout_seconds       = 60
      period_seconds        = 30
      failure_threshold     = 6
      check_method          = "HTTP"
      url                   = "/health"
    }
    liveness_health = {
      initial_delay_seconds = 60
      timeout_seconds       = 60
      period_seconds        = 30
      failure_threshold     = 12
      check_method          = "HTTP"
      url                   = "/health"
    }
    startup_health = {
      initial_delay_seconds = 300
      timeout_seconds       = 60
      period_seconds        = 60
      failure_threshold     = 200
      check_method          = "HTTP"
      url                   = "/health"
    }
  }
]

service_runtime_config = <<-JSON
{
  "service_invoke": {
    "port": 8080,
    "protocol": "HTTPS",
    "auth_type": "TOKEN",
    "direct_channel_auth_enable": false
  },
  "service_limit": {
    "request_size_limit": 20,
    "request_timeout": 30,
    "ip_white_list": [],
    "ip_black_list": [],
    "rate_limit": {
      "num": 200,
      "unit": "SECONDS"
    }
  }
}
JSON

service_upgrade_config = <<-JSON
{
  "type": "ROLLING",
  "rolling_update": {
    "max_surge": "50%",
    "max_unavailable": "50%"
  }
}
JSON

service_log_configs = [
  {
    type = "STDOUT"
  }
]

service_tags = {
  source = "terraform"
}
