terraform {
  required_version = ">= 0.14.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">=1.66.3"
    }
    random = {
      source  = "hashicorp/random"
      version = ">=3.0.0"
    }
  }
}

# Master instance located by
provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}

# Stanby instance located by
provider "huaweicloud" {
  alias      = "dr"
  region     = var.dr_region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
