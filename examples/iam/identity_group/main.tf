#Create two Group
resource "huaweicloud_identity_group_v3" "group" {
  count       = length(var.iden_group)

  name        = lookup(var.iden_group[count.index], "name", null)
  description = lookup(var.iden_group[count.index], "description", null)
}

resource "huaweicloud_identity_group_v3" "group_default" {
  name        = var.iden_group_default_name
  description = var.iden_group_default_description
}