terraform {
  required_providers {
    required_version = ">= 1.3.0"

    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">=1.58.0"
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
