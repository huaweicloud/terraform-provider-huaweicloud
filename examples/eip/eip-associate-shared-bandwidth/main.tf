resource "huaweicloud_vpc_bandwidth" "test" {
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  name                  = var.bandwidth_name
  size                  = var.bandwidth_size
  charge_mode           = var.bandwidth_charge_mode
  bandwidth_type        = var.bandwidth_type
  public_border_group   = var.bandwidth_public_border_group
}

resource "huaweicloud_vpc_eip" "test" {
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.eip_bandwidth_name
    size        = var.eip_bandwidth_size
    share_type  = "PER"
    charge_mode = var.eip_bandwidth_charge_mode
  }

  description = var.eip_description
  tags        = var.eip_tags

  # After adding to the shared bandwidth, `bandwidth.share_type` will be automatically set to `WHOLE`, so the change needs to be ignored.
  lifecycle {
    ignore_changes = [bandwidth]
  }
}

resource "huaweicloud_eip_bandwidth_associate" "test" {
  publicip_id           = huaweicloud_vpc_eip.test.id
  bandwidth_id          = huaweicloud_vpc_bandwidth.test.id
  bandwidth_charge_mode = var.eip_bandwidth_charge_mode
  bandwidth_size        = var.eip_bandwidth_size
  bandwidth_name        = var.eip_bandwidth_name
}
