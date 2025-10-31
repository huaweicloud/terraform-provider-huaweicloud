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
  host                   = "https://%{if var.eip_address != ""}${var.eip_address}:5443%{else}${huaweicloud_vpc_eip.test[0].address}:5443%{endif}"
  cluster_ca_certificate = base64decode(huaweicloud_cce_cluster.test.certificate_clusters[0].certificate_authority_data)
  client_certificate     = base64decode(huaweicloud_cce_cluster.test.certificate_users[0].client_certificate_data)
  client_key             = base64decode(huaweicloud_cce_cluster.test.certificate_users[0].client_key_data)
}
