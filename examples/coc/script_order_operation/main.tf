resource "huaweicloud_coc_script" "test" {
  name        = var.coc_script_name
  description = var.coc_script_description
  risk_level  = var.coc_script_risk_level
  version     = var.coc_script_version
  type        = var.coc_script_type
  content     = var.coc_script_content

  dynamic "parameters" {
    for_each = var.coc_script_parameters

    content {
      name        = parameters.value.name
      value       = parameters.value.value
      description = parameters.value.description
      sensitive   = parameters.value.sensitive
    }
  }
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "")
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_vpc" "test" {
  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "")
  name              = var.subnet_name
  cidr              = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip        = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

data "huaweicloud_images_images" "test" {
  flavor_id  = try(data.huaweicloud_compute_flavors.test.ids[0], "")
  os         = "Ubuntu"
  visibility = "public"
}

# The default security group rules cannot be deleted, otherwise the UniAgent installation will fail.
resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

# Create an ECS instance and install UniAgent
resource "huaweicloud_compute_instance" "test" {
  name                  = var.ecs_instance_name
  availability_zone     = try(data.huaweicloud_availability_zones.test.names[0], "")
  flavor_id             = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  image_id              = try(data.huaweicloud_images_images.test.images[0].id, "")
  security_group_ids    = [huaweicloud_networking_secgroup.test.id]
  user_data             = var.ecs_instance_user_data
  enterprise_project_id = var.enterprise_project_id

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = huaweicloud_compute_instance.test.id
  timeout      = var.coc_script_execute_timeout
  execute_user = var.coc_script_execute_user

  dynamic "parameters" {
    for_each = var.coc_script_execute_parameters

    content {
      name  = parameters.value.name
      value = parameters.value.value
    }
  }

  timeouts {
    create = "1m"
  }
}

resource "huaweicloud_coc_script_order_operation" "test" {
  execute_uuid   = huaweicloud_coc_script_execute.test.id
  batch_index    = var.coc_script_order_operation_batch_index
  instance_id    = var.coc_script_order_operation_instance_id
  operation_type = var.coc_script_order_operation_type
}
