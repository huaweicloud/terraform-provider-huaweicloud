data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = var.organizational_unit_name
  parent_id = data.huaweicloud_organizations_organization.test.root_id
  tags      = var.tags
}
