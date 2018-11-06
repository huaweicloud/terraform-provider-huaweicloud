provider "huaweicloud" {
  user_name   = "${var.user_name}"
  tenant_name = "${var.tenant_name}"
  password    = "${var.password}"
  auth_url    = "${var.auth_url}"
  region      = "${var.region}"
  domain_name = "${var.domain_name}"
}

resource "huaweicloud_blockstorage_volume_v2" "volume" {
  count = "${var.disk_size_gb > 0 ? var.instance_count : 0}"
  name  = "${var.project}-disk${format("%02d", count.index+1)}"
  size  = "${var.disk_size_gb}"
  tags {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_vbs_backup_v2" "backups_1" {
  volume_id = "${huaweicloud_blockstorage_volume_v2.volume.id}"
  name = "${var.project}-backup${format("%02d", count.index+1)}"
}

data "huaweicloud_vbs_backup_v2" "backups" {
  id = "${huaweicloud_vbs_backup_v2.backups_1.id}"
}

resource "huaweicloud_vbs_backup_policy_v2" "vbs" {
  name = "${var.project}-backup-policy${format("%02d", count.index+1)}"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
  tags =[
    {
      key = "k2"
      value = "v2"
    }]
}

data "huaweicloud_vbs_backup_policy_v2" "policies" {
  id = "${huaweicloud_vbs_backup_policy_v2.vbs.id}"
}