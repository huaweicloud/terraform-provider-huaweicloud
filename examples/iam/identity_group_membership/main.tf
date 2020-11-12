#Create a User
resource "huaweicloud_identity_user_v3" "user_A" {
  name        = var.user_A_name
  password    = var.user_A_password
}

#Create two Group
resource "huaweicloud_identity_group_v3" "group" {
  count       = length(var.iden_group)

  name        = lookup(var.iden_group[count.index], "name", null)
  description = lookup(var.iden_group[count.index], "description", null)
}

#Create a default Group
resource "huaweicloud_identity_group_v3" "default_group" {
  name        = "default_iden_group"
  description = "This is a default identity group."
}

#Create MemberShip
resource "huaweicloud_identity_group_membership_v3" "membership_1" {
  group = length(huaweicloud_identity_group_v3.group)==2? huaweicloud_identity_group_v3.group[1].id : huaweicloud_identity_group_v3.default_group.id
  users = [huaweicloud_identity_user_v3.user_A.id]
}
