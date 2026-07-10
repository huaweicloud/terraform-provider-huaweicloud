terraform {
  required_version = ">= 1.9.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # SC.004 Disable
      # 1.80.5 resolves the issue with the status field when unmarshaling the details of the CCE cluster
      version = ">= 1.80.5"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
