# Create an organization
resource "huaweicloud_organizations_organization" "test" {
  enabled_policy_types = var.enabled_policy_types
  root_tags            = var.root_tags
}
