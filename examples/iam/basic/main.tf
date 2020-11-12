
resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "_%@"
}

resource "huaweicloud_identity_user_v3" "iden_user_default" {
  name        = var.iden_user_name_default
  password    = random_password.password.result
}

resource "huaweicloud_identity_user_v3" "users" {
  count       = length(var.iden_users)

  name        = lookup(var.iden_users[count.index], "name", null)
  description = lookup(var.iden_users[count.index], "description", null)
  password    = random_password.password.result
}
