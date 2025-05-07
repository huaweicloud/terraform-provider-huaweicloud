data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_sms_source_servers" "source" {
  name = var.sms_source_server_name
}

resource "huaweicloud_sms_server_template" "test" {
  name              = var.sms_server_template_name
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_sms_task" "migration" {
  type             = "MIGRATE_FILE"
  os_type          = "LINUX"
  source_server_id = data.huaweicloud_sms_source_servers.source.servers[0].id
  vm_template_id   = huaweicloud_sms_server_template.test.id
}
