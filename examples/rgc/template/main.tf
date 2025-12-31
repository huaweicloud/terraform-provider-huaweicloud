resource "huaweicloud_rgc_template" "test" {
  template_name        = var.template_name
  template_type        = var.template_type
  template_description = var.template_description
  template_body        = var.template_body
}
