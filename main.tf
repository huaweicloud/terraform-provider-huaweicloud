terraform{
  required_providers {
    huaweicloud = {
      source = "huaweicloud/huaweicloud"
    }
  }
}

provider "huaweicloud"{
  region = "cn-north-4"
  access_key = var.access_key
  secret_key = var.secret_key
}

data "huaweicloud_dsc_security_levels" "test" {}