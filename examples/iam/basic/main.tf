data "huaweicloud_identity_role" "auth_admin" {
  name = "system_all_0"
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "_%@"
}

resource "huaweicloud_identity_user" "user_A" {
  name     = var.iden_user_name
  password = random_password.password.result
}

resource "huaweicloud_identity_group" "group" {
  count = length(var.iden_group)

  name        = lookup(var.iden_group[count.index], "name", null)
  description = lookup(var.iden_group[count.index], "description", null)
}

resource "huaweicloud_identity_group" "default_group" {
  name        = "default_group"
  description = "This is a default identity group."
}

resource "huaweicloud_identity_group_membership" "membership_1" {
  group = length(huaweicloud_identity_group.group) >= 2 ? huaweicloud_identity_group.group[1].id : huaweicloud_identity_group.default_group.id
  users = [huaweicloud_identity_user.user_A.id]
}

resource "huaweicloud_identity_role_assignment" "role_assignment_1" {
  group_id  = length(huaweicloud_identity_group.group) >= 2 ? huaweicloud_identity_group.group[1].id : huaweicloud_identity_group.default_group.id
  domain_id = var.domain_id
  role_id   = data.huaweicloud_identity_role.auth_admin.id
}
