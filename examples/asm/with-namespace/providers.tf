terraform {
  required_version = ">= 1.9.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # SC.004 Disable
      # In 1.74.1, the API will return an invalid response format and cause timeout
      version = "= 1.75.0"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
