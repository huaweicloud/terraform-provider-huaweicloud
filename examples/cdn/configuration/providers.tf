terraform {
  required_providers {
    huaweicloud = {
      source  = "local-registry/huaweicloud/huaweicloud"
      version = ">=1.70.1"
    }
  }
}

provider "huaweicloud" {
  region     = "cn-north-4"
  access_key = var.access_key
  secret_key = var.secret_key
}
