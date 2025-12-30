resource "huaweicloud_vpn_user" "test" {
  vpn_server_id = var.vpn_user_server_id
  name          = var.vpn_user_name
  password      = var.vpn_user_password
}
