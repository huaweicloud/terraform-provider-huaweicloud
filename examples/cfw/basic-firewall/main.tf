# Purchase a Cloud Firewall instance
resource "huaweicloud_cfw_firewall" "test" {
  name = var.firewall_name

  flavor {
    version = var.firewall_flavor
  }

  charging_mode = var.firewall_charging_mode
  tags          = var.firewall_tags
}

# Enable EIP auto-protection
resource "huaweicloud_cfw_eip_auto_protection" "test" {
  fw_instance_id = huaweicloud_cfw_firewall.test.id
  object_id      = try(huaweicloud_cfw_firewall.test.protect_objects[0].object_id, "")
  status         = var.eip_auto_protection_status
}

# Manually bind existing EIPs for protection
resource "huaweicloud_cfw_eip_protection" "test" {
  count = var.eip_protection_enabled && length(var.eip_protection_eip_ids) > 0 ? 1 : 0

  object_id = try(huaweicloud_cfw_firewall.test.protect_objects[0].object_id, "")

  dynamic "protected_eip" {
    for_each = var.eip_protection_eip_ids

    content {
      id          = protected_eip.value.id
      public_ipv4 = protected_eip.value.public_ipv4
    }
  }
}
