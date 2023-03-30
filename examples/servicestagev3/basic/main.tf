resource "huaweicloud_servicestage_v3_application" "application_terraform_1" {
  name                  = var.app_name
  description           = var.app_description
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_servicestage_v3_environment" "environment_terraform_1" {
  name                  = var.env_name
  description           = var.env_description
  deploy_mode           = var.env_deploy_mode
  enterprise_project_id = var.enterprise_project_id
  vpc_id                = var.vpc_id

  dynamic "labels" {
    for_each = var.env_labels

    content {
      key   = labels.value.key
      value = labels.value.value
    }
  }
}