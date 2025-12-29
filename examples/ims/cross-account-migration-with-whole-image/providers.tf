terraform {
  required_version = ">= 1.1.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">= 1.68.0"
    }
  }
}

# Provider for sharer (source account)
provider "huaweicloud" {
  alias      = "sharer"
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}

# Provider for accepter (target account)
provider "huaweicloud" {
  alias      = "accepter"
  region     = var.region_name
  access_key = var.accepter_access_key
  secret_key = var.accepter_secret_key
}
