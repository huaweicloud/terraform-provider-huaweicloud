terraform {
  required_version = ">= 1.3.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # When the version is lower than 1.94.0, the huaweicloud_modelartsv2_service resource will report an error when
      # using 'log_configs' parameter.
      # SC.004 Disable
      version = ">= 1.94.0"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
