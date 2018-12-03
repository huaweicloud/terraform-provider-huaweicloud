resource "huaweicloud_networking_network_v2" "terraform" {
  name           = "terraform"
  admin_state_up = "true"
}

resource "huaweicloud_compute_secgroup_v2" "terraform" {
  name        = "terraform"
  description = "Security group for the Terraform example instances"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}

resource "huaweicloud_compute_instance_v2" "terraform" {
  name            = "terraform"
  image_name      = "${var.image}"
  flavor_name     = "${var.flavor}"
  availability_zone = "${var.availability_zone}"
  security_groups = ["${huaweicloud_compute_secgroup_v2.terraform.name}"]

  network {
    uuid = "${huaweicloud_networking_network_v2.terraform.id}"
  }

}
resource "huaweicloud_csbs_backup_v1" "backup_v1" {
  backup_name      = "${var.project}-backup"
  resource_id      = "${huaweicloud_compute_instance_v2.terraform.id}"
  resource_type    = "OS::Nova::Server"
}

resource "huaweicloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name            = "csbs-backup-policy"
  resource {
    id = "${huaweicloud_compute_instance_v2.terraform.id}"
    type = "OS::Nova::Server"
    name = "resource1"
  }
  scheduled_operation {
    name ="mybackup"
    enabled = true
    operation_type ="backup"
    max_backups = 2
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}