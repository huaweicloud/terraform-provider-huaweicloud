terraform {
  required_version = ">= 0.14.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      # When the version is lower than 1.80.1, the huaweicloud_dns_zone_authorization resource will throw an error
      # during the query phase. The higher version make this resource to a one-time action resource to avoid this error
      # reported (Only parse the response of the request in the CreateContext phase).
      # SC.004 Disable
      version = ">= 1.80.1"
      # SC.004 Enable
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}

provider "huaweicloud" {
  alias = "domain_master"

  region     = var.region_name
  access_key = var.target_account_access_key
  secret_key = var.target_account_secret_key
}
