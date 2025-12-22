# Create a SWR repository
resource "huaweicloud_swr_organization" "test" {
  name = var.organization_name
}

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = var.repository_name
  category     = var.category
}
