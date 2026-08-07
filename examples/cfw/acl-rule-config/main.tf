# Query existing firewall information
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = var.fw_instance_id != "" ? var.fw_instance_id : null
}

locals {
  fw_instance_id = var.fw_instance_id != "" ? var.fw_instance_id : try(data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id, null)
  object_id      = try(data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id, null)
}

# Create an IP address group
resource "huaweicloud_cfw_address_group" "test" {
  object_id   = local.object_id
  name        = var.address_group_name
  description = var.address_group_description
}

# Create a service group
resource "huaweicloud_cfw_service_group" "test" {
  object_id   = local.object_id
  name        = var.service_group_name
  description = var.service_group_description
}

# Create a domain name group
resource "huaweicloud_cfw_domain_name_group" "test" {
  fw_instance_id = local.fw_instance_id
  object_id      = local.object_id
  name           = var.domain_name_group_name
  type           = var.domain_name_group_type
  description    = var.domain_name_group_description

  dynamic "domain_names" {
    for_each = var.domain_name_group_domains

    content {
      domain_name = domain_names.value.domain_name
      description = domain_names.value.description
    }
  }
}

# ST.001 Disable
# Create an ACL rule with IP-based source and destination
resource "huaweicloud_cfw_acl_rule" "ip_based" {
  name                = var.acl_rule_ip_name
  object_id           = local.object_id
  description         = var.acl_rule_ip_description
  type                = var.acl_rule_type
  address_type        = var.acl_rule_address_type
  action_type         = var.acl_rule_action_type
  long_connect_enable = var.acl_rule_long_connect_enable
  status              = var.acl_rule_status
  applications        = var.acl_rule_applications

  source_addresses      = var.acl_rule_source_addresses
  destination_addresses = var.acl_rule_destination_addresses

  custom_services {
    protocol    = var.acl_rule_custom_service_protocol
    source_port = var.acl_rule_custom_service_source_port
    dest_port   = var.acl_rule_custom_service_dest_port
  }

  sequence {
    top = 1
  }

  tags = var.tags
}

# Create an ACL rule with domain-based destination
resource "huaweicloud_cfw_acl_rule" "domain_based" {
  name                = var.acl_rule_domain_name
  object_id           = local.object_id
  description         = var.acl_rule_domain_description
  type                = var.acl_rule_type
  address_type        = var.acl_rule_address_type
  action_type         = var.acl_rule_action_type
  long_connect_enable = var.acl_rule_long_connect_enable
  status              = var.acl_rule_status
  direction           = var.acl_rule_domain_direction

  source_addresses                = var.acl_rule_source_addresses
  destination_domain_address_name = var.acl_rule_destination_domain_address_name

  custom_services {
    protocol    = var.acl_rule_custom_service_protocol
    source_port = var.acl_rule_custom_service_source_port
    dest_port   = var.acl_rule_custom_service_dest_port
  }

  sequence {
    top          = 0
    dest_rule_id = huaweicloud_cfw_acl_rule.ip_based.id
  }

  tags = var.tags
}

# Create an ACL rule with address group-based source and destination
resource "huaweicloud_cfw_acl_rule" "group_based" {
  name                = var.acl_rule_group_name
  object_id           = local.object_id
  description         = var.acl_rule_group_description
  type                = var.acl_rule_type
  address_type        = var.acl_rule_address_type
  action_type         = var.acl_rule_action_type
  long_connect_enable = var.acl_rule_long_connect_enable
  status              = var.acl_rule_status

  source_address_groups      = [huaweicloud_cfw_address_group.test.id]
  destination_address_groups = [huaweicloud_cfw_address_group.test.id]

  custom_service_groups {
    protocols = [var.acl_rule_service_group_protocol]
    group_ids = [huaweicloud_cfw_service_group.test.id]
  }

  sequence {
    bottom = 1
  }

  tags = var.tags
}
# ST.001 Enable
