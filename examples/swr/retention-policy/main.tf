# Create a SWR repository image retention policy
resource "huaweicloud_swr_organization" "test" {
  name = var.organization_name
}

resource "huaweicloud_swr_repository" "test" {
  organization = huaweicloud_swr_organization.test.name
  name         = var.repository_name
  category     = var.category
}

resource "huaweicloud_swr_image_retention_policy" "test" {
  organization = huaweicloud_swr_organization.test.name
  repository   = huaweicloud_swr_repository.test.name
  type         = var.policy_type
  number       = var.policy_number

  dynamic "tag_selectors" {
    for_each = var.tag_selectors_configuration

    content {
      kind    = tag_selectors.value.kind
      pattern = tag_selectors.value.pattern
    }
  }
}
