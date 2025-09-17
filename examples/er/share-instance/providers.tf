terraform {
  required_version = ">= 1.1.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">= 1.73.4"
    }
  }
}

# Share (owner).
provider "huaweicloud" {
  alias = "owner"

  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}

# Other account (principal).
provider "huaweicloud" {
  alias = "principal"

  region     = var.region_name
  access_key = var.principal_access_key
  secret_key = var.principal_secret_key
}
