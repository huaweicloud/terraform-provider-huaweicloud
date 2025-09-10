terraform {
  required_providers {
    required_version = ">= 1.1.0"

    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">= 1.68.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.13"
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}
