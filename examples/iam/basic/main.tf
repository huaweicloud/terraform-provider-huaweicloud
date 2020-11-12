#Create a User
resource "huaweicloud_identity_user_v3" "iden_user_default" {
  name        = var.iden_user_name_default
  password    = var.iden_user_password_default
}

#Create two Users
resource "huaweicloud_identity_user_v3" "users" {
  count       = length(var.iden_users)

  name        = lookup(var.iden_users[count.index], "name", null)
  description = lookup(var.iden_users[count.index], "description", null)
  password    = lookup(var.iden_users[count.index], "password", null)
}
