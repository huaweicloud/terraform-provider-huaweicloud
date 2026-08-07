# Query existing firewall information
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = var.fw_instance_id != "" ? var.fw_instance_id : null
}

locals {
  object_id = try(data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id, null)
}

# Create a blacklist rule to block specific IP addresses
resource "huaweicloud_cfw_black_white_list" "test" {
  count = 2

  object_id    = local.object_id
  list_type    = count.index == 0 ? var.blacklist_list_type : var.whitelist_list_type
  direction    = count.index == 0 ? var.blacklist_direction : var.whitelist_direction
  protocol     = count.index == 0 ? var.blacklist_protocol  : var.whitelist_protocol
  port         = count.index == 0 ? var.blacklist_port : var.whitelist_port
  address_type = count.index == 0 ? var.blacklist_address_type : var.whitelist_address_type
  address      = count.index == 0 ? var.blacklist_address : var.whitelist_address
}
