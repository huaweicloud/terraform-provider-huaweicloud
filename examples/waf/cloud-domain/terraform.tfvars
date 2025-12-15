cloud_instance_resource_spec_code = "detection"
cloud_instance_charging_mode      = "prePaid"
cloud_instance_period_unit        = "month"
cloud_instance_period             = 1
cloud_domain                      = "demo-example-test.huawei.com"

cloud_server = [
  {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "119.8.0.17"
    port            = "8080"
    type            = "ipv4"
    weight          = 1
  }
]
