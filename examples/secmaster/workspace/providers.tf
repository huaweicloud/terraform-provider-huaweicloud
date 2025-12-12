terraform {
  required_version = ">= 0.14.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # When the version is lower than 1.76.0, the huaweicloud_secmaster_workspace resource is a one-time action
      # resource, execute `terraform destroy` unable to delete the workspace.The higher version make this resource
      # to a normal (Create, Read, Update, Delete) resource to avoid this issue.
      # SC.004 Disable
      version = ">= 1.76.0"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
