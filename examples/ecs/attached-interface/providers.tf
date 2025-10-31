terraform {
  required_version = ">= 1.9.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # When the version is lower than 1.61.0, the huaweicloud_compute_interface_attach resource will throw an error during the creation phase.
      # SC.004 Disable
      version = ">= 1.61.0"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
