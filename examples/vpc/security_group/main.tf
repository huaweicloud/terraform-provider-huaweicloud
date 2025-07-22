resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = length(var.security_group_rule_configurations)

  direction         = lookup(var.security_group_rule_configurations[count.index], "direction", "ingress")
  ethertype         = lookup(var.security_group_rule_configurations[count.index], "ethertype", "IPv4")
  protocol          = lookup(var.security_group_rule_configurations[count.index], "protocol", null)
  ports             = lookup(var.security_group_rule_configurations[count.index], "ports", null)
  remote_ip_prefix  = lookup(var.security_group_rule_configurations[count.index], "remote_ip_prefix", "0.0.0.0/0")
  security_group_id = huaweicloud_networking_secgroup.test.id
}
