terraform {
  required_version = ">= 1.9.0"

  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">= 1.57.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 1.6.2"
    }
  }
}

provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}

provider "kubernetes" {
  config_path    = local_file.test.filename
  config_context = "external"
}
