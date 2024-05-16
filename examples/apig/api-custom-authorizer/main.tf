data "huaweicloud_availability_zones" "default" {}

resource "huaweicloud_vpc" "default" {
  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "default" {
  name       = var.subnet_name
  cidr       = "192.168.0.0/20"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.default.id
}

resource "huaweicloud_networking_secgroup" "default" {
  name = var.security_group_name
}

# HTTP access from anywhere
resource "huaweicloud_networking_secgroup_rule" "http_rule" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_max    = 80
  port_range_min    = 80
  remote_ip_prefix  = "0.0.0.0/0"
}

data "huaweicloud_images_image" "default" {
  name        = var.ecs_image
  most_recent = true
}

data "huaweicloud_compute_flavors" "default" {
  availability_zone = data.huaweicloud_availability_zones.default.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "!@"
  min_numeric      = 1
  min_lower        = 1
  min_special      = 1
}

resource "huaweicloud_compute_instance" "default" {
  name               = var.ecs_name
  image_id           = data.huaweicloud_images_image.default.id
  flavor_id          = data.huaweicloud_compute_flavors.default.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.default.id]
  availability_zone  = data.huaweicloud_availability_zones.default.names[0]
  system_disk_type   = "SSD"
  admin_pass         = random_password.password.result

  network {
    uuid = huaweicloud_vpc_subnet.default.id
  }
}

resource "huaweicloud_vpc_eip" "default" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = var.bandwidth_name
    size        = 10
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_eip_associate" "default" {
  public_ip   = huaweicloud_vpc_eip.default.address
  instance_id = huaweicloud_compute_instance.default.id
}

resource "null_resource" "provision" {
  depends_on = [
    huaweicloud_compute_eip_associate.default
  ]
  provisioner "remote-exec" {
    connection {
      user     = "root"
      password = random_password.password.result
      host     = huaweicloud_vpc_eip.default.address
      port     = 22
    }
    inline = [
      "yum -y install nginx",
      "systemctl enable nginx",
      "systemctl start nginx",
      "systemctl status nginx",
    ]
  }
}

resource "huaweicloud_fgs_function" "default" {
  name        = var.function_name
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python3.6"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler(event, context):
    if event["headers"].get("x-user-auth")=='cXpsdzQyVW9Xa1NVTX==':
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"allow",
                "context":{
                    "user_name":"user1"
                }
            })
        }
    else:
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"deny",
                "context":{
                    "code":"1001",
                    "message":"incorrect username or password"
                }
            })
        }
EOF
}

resource "huaweicloud_apig_instance" "default" {
  name              = var.apig_instance_name
  edition           = "BASIC"
  vpc_id            = huaweicloud_vpc.default.id
  subnet_id         = huaweicloud_vpc_subnet.default.id
  security_group_id = huaweicloud_networking_secgroup.default.id

  available_zones = [
    data.huaweicloud_availability_zones.default.names[0],
  ]
}

resource "huaweicloud_apig_custom_authorizer" "default" {
  instance_id      = huaweicloud_apig_instance.default.id
  name             = var.apig_auth_name
  function_urn     = huaweicloud_fgs_function.default.urn
  function_version = "latest"
  type             = "FRONTEND"
}

resource "huaweicloud_apig_response" "default" {
  name        = var.apig_response_name
  instance_id = huaweicloud_apig_instance.default.id
  group_id    = huaweicloud_apig_group.default.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}

resource "huaweicloud_apig_group" "default" {
  name        = var.apig_group_name
  instance_id = huaweicloud_apig_instance.default.id
}

resource "huaweicloud_apig_vpc_channel" "default" {
  instance_id = huaweicloud_apig_instance.default.id
  name        = var.apig_channel_name
  port        = 80
  member_type = "ECS"
  protocol    = "HTTP"
  path        = "/backend/users"
  http_code   = "200,401"

  members {
    id = huaweicloud_compute_instance.default.id
  }
}

resource "huaweicloud_apig_api" "default" {
  instance_id             = huaweicloud_apig_instance.default.id
  group_id                = huaweicloud_apig_group.default.id
  type                    = "Public"
  name                    = var.apig_api_name
  request_protocol        = "BOTH"
  request_method          = "GET"
  request_path            = "/terraform/users"
  security_authentication = "AUTHORIZER"
  matching                = "Exact"
  response_id             = huaweicloud_apig_response.default.id
  authorizer_id           = huaweicloud_apig_custom_authorizer.default.id

  backend_params {
    type     = "SYSTEM"
    name     = "X-User-Auth"
    location = "HEADER"
    value    = "user_name"
  }

  web {
    # ensure the backend server include this path or API debug will fail
    path             = "/backend/users"
    vpc_channel_id   = huaweicloud_apig_vpc_channel.default.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 5000
  }
}
