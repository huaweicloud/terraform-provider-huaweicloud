terraform {
  required_version = ">= 1.1.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # SC.004 Disable
      # SC.004 Error: test case in local can passed, but in CI will failed.
      version = ">= 1.69.0"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
