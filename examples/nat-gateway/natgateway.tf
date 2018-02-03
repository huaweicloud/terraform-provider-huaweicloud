resource "huaweicloud_nat_gateway_v2" "nat_1" {
  name   = "Terraform"
  description = "test for terraform"
  spec = "1"
  router_id = "2c1fe4bd-ebad-44ca-ae9d-e94e63847b75"
  internal_network_id = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
}

resource "huaweicloud_nat_snat_rule_v2" "rule_1" {
  nat_gateway_id = "${huaweicloud_vpc_nat_gateway_v2.nat_1.id}"
  network_id = "dc8632e2-d9ff-41b1-aa0c-d455557314a0"
  floating_ip_id = "0a166fc5-a904-42fb-b1ef-cf18afeeddca"
}
