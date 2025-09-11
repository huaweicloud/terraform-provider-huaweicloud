# Create VPC for testing policy assignment
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

# Create subnet for VPC
resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

# Create security group
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

# Create ECS instance for policy assignment testing
resource "huaweicloud_compute_instance" "test" {
  name               = var.ecs_instance_name
  image_name         = var.ecs_image_name
  flavor_name        = var.ecs_flavor_name
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = var.availability_zone

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  tags = var.ecs_tags
}

# Create a policy assignment to check if ECS instances have specific tags
resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = var.policy_assignment_name
  description          = var.policy_assignment_description
  policy_definition_id = var.policy_definition_id

  dynamic "policy_filter" {
    for_each = var.policy_assignment_policy_filter

    content {
      region            = policy_filter.value.region
      resource_provider = policy_filter.value.resource_provider
      resource_type     = policy_filter.value.resource_type
      resource_id       = policy_filter.value.resource_id
      tag_key           = policy_filter.value.tag_key
      tag_value         = policy_filter.value.tag_value
    }
  }

  parameters = var.policy_assignment_parameters

  tags = var.policy_assignment_tags
}

# Evaluate the policy assignment
resource "huaweicloud_rms_policy_assignment_evaluate" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
}
