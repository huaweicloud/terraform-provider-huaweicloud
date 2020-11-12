#Automatic password generation
resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "_%@"
}

resource "huaweicloud_identity_user_v3" "user_A" {
  name        = var.iden_user_name
  password    = random_password.password.result
}

resource "huaweicloud_identity_group_v3" "group" {
  count       = length(var.iden_group)

  name        = lookup(var.iden_group[count.index], "name", null)
  description = lookup(var.iden_group[count.index], "description", null)
}

#Create a default Group
resource "huaweicloud_identity_group_v3" "default_group" {
  name        = "default_group"
  description = "This is a default identity group."
}

#  Add user_A to second identity group
#  When second custom identity groups do not exist, add user_A to the default identity group
resource "huaweicloud_identity_group_membership_v3" "membership_1" {
  group = length(huaweicloud_identity_group_v3.group)>=2? huaweicloud_identity_group_v3.group[1].id : huaweicloud_identity_group_v3.default_group.id
  users = [huaweicloud_identity_user_v3.user_A.id]
}

data "huaweicloud_identity_role_v3" "auth_admin" {
    name = "system_all_0"
}

resource "huaweicloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id  = length(huaweicloud_identity_group_v3.group)>=2? huaweicloud_identity_group_v3.group[1].id : huaweicloud_identity_group_v3.default_group.id
  domain_id = var.domain_id
  role_id   = data.huaweicloud_identity_role_v3.auth_admin.id
}
