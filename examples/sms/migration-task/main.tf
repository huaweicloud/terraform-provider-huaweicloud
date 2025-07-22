data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_sms_source_servers" "test" {
  name = var.source_server_name
}

resource "huaweicloud_sms_server_template" "test" {
  name              = var.server_template_name
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
}

resource "huaweicloud_sms_task" "test" {
  type             = var.migrate_task_type
  os_type          = var.server_os_type
  source_server_id = try(data.huaweicloud_sms_source_servers.test.servers[0].id, null)
  vm_template_id   = huaweicloud_sms_server_template.test.id
}
