data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

resource "huaweicloud_vpc" "test" {
  count = var.vpc_id == "" && var.subnet_id == "" ? 1 : 0

  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = var.subnet_id == "" ? 1 : 0

  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  name              = var.subnet_name
  cidr              = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0)
  gateway_ip        = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : var.subnet_cidr != "" ? cidrhost(var.subnet_cidr, 1) : cidrhost(cidrsubnet(huaweicloud_vpc.test[0].cidr, 4, 0), 1)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_vpc_eip" "test" {
  count = var.eip_address == "" ? 1 : 0

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = var.cluster_name
  flavor_id              = var.cluster_flavor_id
  cluster_version        = var.cluster_version
  cluster_type           = var.cluster_type
  container_network_type = var.container_network_type
  vpc_id                 = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  subnet_id              = var.subnet_id != "" ? var.subnet_id : huaweicloud_vpc_subnet.test[0].id
  eip                    = var.eip_address != "" ? var.eip_address : huaweicloud_vpc_eip.test[0].address
  authentication_mode    = var.authentication_mode
  delete_all             = var.delete_all_resources_on_terminal ? "true" : "false"
  enterprise_project_id  = var.enterprise_project_id
}

data "huaweicloud_compute_flavors" "test" {
  count = var.node_flavor_id == "" ? 1 : 0

  performance_type  = var.node_performance_type
  cpu_core_count    = var.node_cpu_core_count
  memory_size       = var.node_memory_size
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_kps_keypair" "test" {
  name = var.keypair_name
}

resource "huaweicloud_cce_node" "test" {
  cluster_id        = huaweicloud_cce_cluster.test.id
  name              = var.node_name
  flavor_id         = var.node_flavor_id != "" ? var.node_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  key_pair          = huaweicloud_kps_keypair.test.name

  root_volume {
    volumetype = var.root_volume_type
    size       = var.root_volume_size
  }

  dynamic "data_volumes" {
    for_each = var.data_volumes_configuration

    content {
      volumetype = data_volumes.value.volumetype
      size       = data_volumes.value.size
    }
  }
}

resource "kubernetes_secret" "test" {
  metadata {
    name      = var.secret_name
    namespace = var.namespace_name
    labels    = var.secret_labels
  }

  data = var.secret_data
  type = var.secret_type

  # For `data` parameters Input value not consistent with the state value, so we need to ignore the changes.
  lifecycle {
    ignore_changes = [data]
  }
}

resource "kubernetes_persistent_volume_claim" "test" {
  metadata {
    name      = var.pvc_name
    namespace = var.namespace_name

    annotations = {
      "everest.io/obs-volume-type"                       = var.pvc_obs_volume_type
      "csi.storage.k8s.io/fstype"                        = var.pvc_fstype
      "csi.storage.k8s.io/node-publish-secret-name"      = kubernetes_secret.test.metadata[0].name
      "csi.storage.k8s.io/node-publish-secret-namespace" = var.namespace_name
      "everest.io/enterprise-project-id"                 = var.enterprise_project_id
    }
  }

  spec {
    access_modes = var.pvc_access_modes

    resources {
      requests = {
        storage = var.pvc_storage
      }
    }

    storage_class_name = var.pvc_storage_class_name
  }
}

resource "kubernetes_deployment" "test" {
  metadata {
    name      = var.deployment_name
    namespace = var.namespace_name
  }

  spec {
    replicas = var.deployment_replicas

    selector {
      match_labels = {
        app = var.deployment_name
      }
    }

    template {
      metadata {
        labels = {
          app = var.deployment_name
        }
      }

      spec {
        dynamic "container" {
          for_each = var.deployment_containers

          content {
            name  = container.value.name
            image = container.value.image

            dynamic "volume_mount" {
              for_each = container.value.volume_mounts

              content {
                name       = var.deployment_volume_name
                mount_path = volume_mount.value.mount_path
              }
            }
          }
        }

        dynamic "image_pull_secrets" {
          for_each = var.deployment_image_pull_secrets

          content {
            name = image_pull_secrets.value
          }
        }

        volume {
          name = var.deployment_volume_name

          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.test.metadata[0].name
          }
        }
      }
    }
  }

  depends_on = [huaweicloud_cce_node.test]
}
