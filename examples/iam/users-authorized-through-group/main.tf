data "huaweicloud_identity_role" "test" {
  count = var.role_id == "" && var.role_policy == "" ? 1 : 0

  name = var.role_name
}

resource "huaweicloud_identity_role" "test" {
  count = var.role_id == "" && var.role_policy != "" ? 1 : 0

  name        = var.role_name
  type        = var.role_type
  description = var.role_description
  policy      = var.role_policy
}

resource "huaweicloud_identity_group" "test" {
  count = var.group_id == "" ? 1 : 0

  name        = var.group_name
  description = var.group_description
}

data "huaweicloud_identity_projects" "test" {
  count = var.authorized_project_id == "" ? 1 : 0

  name = var.authorized_project_name
}

resource "huaweicloud_identity_role_assignment" "test" {
  group_id   = var.group_id != "" ? var.group_id : huaweicloud_identity_group.test[0].id
  role_id    = var.role_id != "" ? var.role_id : var.role_policy != "" ? huaweicloud_identity_role.test[0].id : try(data.huaweicloud_identity_role.test[0].id, "NOT_FOUND")
  domain_id  = var.authorized_domain_id != "" ? var.authorized_domain_id : null
  project_id = var.authorized_domain_id == "" ? var.authorized_project_id != "" ? var.authorized_project_id : try(data.huaweicloud_identity_projects.test[0].projects[0].id, "NOT_FOUND") : null
}

resource "random_password" "test" {
  count = length([for v in var.users_configuration : true if v.password == "" || v.password == null]) > 0 ? 1 : 0

  length           = 16
  special          = true
  override_special = "_%@"
}

resource "huaweicloud_identity_user" "test" {
  count = length(var.users_configuration)

  name     = lookup(var.users_configuration[count.index], "name", null)
  password = lookup(var.users_configuration[count.index], "password", null) != "" ? lookup(var.users_configuration[count.index], "password", null) : random_password.test[count.index].result
}

resource "huaweicloud_identity_group_membership" "test" {
  group = var.group_id != "" ? var.group_id : huaweicloud_identity_group.test[0].id
  users = huaweicloud_identity_user.test[*].id
}
