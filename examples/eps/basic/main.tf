resource "huaweicloud_enterprise_project" "test" {
  name                    = var.enterprise_project_name
  description             = var.enterprise_project_description
  type                    = var.enterprise_project_type
  enable                  = var.enterprise_project_enable
  skip_disable_on_destroy = var.skip_disable_on_destroy
  delete_flag             = var.delete_flag
}
