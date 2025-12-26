# terraform {
#   required_version = ">= 0.14.0"

#   required_providers {
#     huaweicloud = {
#       source  = "huaweicloud/huaweicloud"
#       version = ">= 1.37.0"
#     }
#   }
# }

terraform {
  required_providers {
    huaweicloud = {
      source  = "local-registry/huaweicloud/huaweicloud"
      version = "0.0.1"
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
