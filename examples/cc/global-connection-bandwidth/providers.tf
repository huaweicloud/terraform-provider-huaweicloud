terraform {
  required_version = ">= 0.14.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # SC.004 Disable
      version = ">= 1.77.2"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
