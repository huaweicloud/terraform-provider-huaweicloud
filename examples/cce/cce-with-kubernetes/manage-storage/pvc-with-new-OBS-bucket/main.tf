data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_vpc" "myvpc" {
  id = var.vpc_id
}

data "huaweicloud_vpc_subnet" "mysubnet" {
  id = var.subnet_id
}

resource "huaweicloud_vpc_eip" "cce" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "cce-apiserver"
    size        = 20
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_cluster" "cluster" {
  name                   = var.cce_name
  cluster_type           = "VirtualMachine"
  cluster_version        = "v1.19"
  flavor_id              = "cce.s1.small"
  vpc_id                 = data.huaweicloud_vpc.myvpc.id
  subnet_id              = data.huaweicloud_vpc_subnet.mysubnet.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.cce.address
  delete_all             = "true"
}

resource "huaweicloud_cce_node" "cce-node1" {
  cluster_id        = huaweicloud_cce_cluster.cluster.id
  name              = "node1"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

resource "huaweicloud_cce_node" "cce-node2" {
  cluster_id        = huaweicloud_cce_cluster.cluster.id
  name              = "node2"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

provider "kubernetes" {
  host                   = "https://${huaweicloud_vpc_eip.cce.address}:5443"
  cluster_ca_certificate = base64decode(huaweicloud_cce_cluster.cluster.certificate_clusters[0].certificate_authority_data)
  client_certificate     = base64decode(huaweicloud_cce_cluster.cluster.certificate_users[0].client_certificate_data)
  client_key             = base64decode(huaweicloud_cce_cluster.cluster.certificate_users[0].client_key_data)
}

resource "kubernetes_secret" "my-secret" {
  metadata {
    name      = "my-secret"
    namespace = "default"

    labels = {
      "secret.kubernetes.io/used-by" = "csi"
    }
  }

  data = {
    "access.key" = "my_access_key"
    "secret.key" = "my_secret_key"
  }

  type = "cfe/secure-opaque"
}

resource "kubernetes_persistent_volume_claim" "my-pvc" {
  metadata {
    name      = "my-pvc-obs"
    namespace = "default"

    annotations = {
      "everest.io/obs-volume-type"                       = "STANDARD"
      "csi.storage.k8s.io/fstype"                        = "s3fs"
      "csi.storage.k8s.io/node-publish-secret-name"      = kubernetes_secret.my-secret.metadata[0].name
      "csi.storage.k8s.io/node-publish-secret-namespace" = kubernetes_secret.my-secret.metadata[0].namespace
      "everest.io/enterprise-project-id"                 = "0"
    }
  }
  spec {
    access_modes = ["ReadWriteMany"]
    resources {
      requests = {
        storage = "1Gi"
      }
    }
    storage_class_name = "csi-obs"
  }
}

resource "kubernetes_deployment" "my-deployment" {
  metadata {
    name      = "web-demo"
    namespace = "default"
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "web-demo"
      }
    }

    template {
      metadata {
        labels = {
          app = "web-demo"
        }
      }

      spec {
        container {
          name  = "container-1"
          image = "nginx:latest"

          volume_mount {
            name       = "pvc-obs-volume"
            mount_path = "/data"
          }
        }
        image_pull_secrets {
          name = "default-secret"
        }
        volume {
          name = "pvc-obs-volume"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.my-pvc.metadata[0].name
          }
        }
      }
    }
  }
}
