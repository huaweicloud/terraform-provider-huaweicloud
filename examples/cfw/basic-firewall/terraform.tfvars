firewall_name              = "tf_test_cfw_firewall"
firewall_flavor            = "Professional"
firewall_charging_mode     = "postPaid"
firewall_tags              = {
  environment = "test"
  managed_by  = "terraform"
}
eip_auto_protection_status = 1
eip_protection_enabled     = false
eip_protection_eip_ids     = []
