data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_deh_types" "test" {
  count = var.instance_host_type == "" ? 1 : 0

  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_deh_instance" "test" {
  name                  = var.instance_name
  availability_zone     = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  host_type             = var.instance_host_type != "" ? var.instance_host_type : try(data.huaweicloud_deh_types.test[0].dedicated_host_types[0].host_type, null)
  auto_placement        = var.instance_auto_placement
  metadata              = var.instance_metadata
  tags                  = var.instance_tags
  enterprise_project_id = var.enterprise_project_id

  charging_mode = var.instance_charging_mode
  period_unit   = var.instance_period_unit
  period        = var.instance_period
  auto_renew    = var.instance_auto_renew
}
