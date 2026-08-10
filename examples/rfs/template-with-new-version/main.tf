resource "huaweicloud_rfs_template" "test" {
  template_name        = var.template_name
  template_body        = var.template_body
  template_description = var.template_description
  version_description  = var.template_initial_version_description
}

resource "huaweicloud_rfs_template_version" "test" {
  template_name       = huaweicloud_rfs_template.test.template_name
  template_id         = huaweicloud_rfs_template.test.template_id
  template_body       = var.template_version_body
  version_description = var.template_version_description
}
