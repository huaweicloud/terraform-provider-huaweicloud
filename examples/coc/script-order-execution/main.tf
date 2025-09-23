data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].ids[0], "") : var.instance_flavor_id
  os         = var.instance_image_os_type
  visibility = var.instance_image_visibility
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
}

# The default security group rules cannot be deleted, otherwise the UniAgent installation will fail.
resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

# Create an ECS instance and install UniAgent
resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  availability_zone  = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, "") : var.instance_flavor_id
  image_id           = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, "") : var.instance_image_id
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  user_data          = var.instance_user_data

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_coc_script" "test" {
  name        = var.script_name
  description = var.script_description
  risk_level  = var.script_risk_level
  version     = var.script_version
  type        = var.script_type
  content     = var.script_content

  dynamic "parameters" {
    for_each = var.script_parameters

    content {
      name        = parameters.value.name
      value       = parameters.value.value
      description = parameters.value.description
      sensitive   = parameters.value.sensitive
    }
  }
}

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = huaweicloud_compute_instance.test.id
  timeout      = var.script_execute_timeout
  execute_user = var.script_execute_user

  dynamic "parameters" {
    for_each = var.script_execute_parameters

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
  batch_index    = var.script_order_operation_batch_index
  instance_id    = var.script_order_operation_instance_id
  operation_type = var.script_order_operation_type
}
