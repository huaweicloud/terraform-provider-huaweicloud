# Create a BCS instance with basic configuration.
resource "huaweicloud_bcs_instance" "test" {
  name                  = var.instance_name
  edition               = var.edition
  fabric_version        = var.fabric_version
  consensus             = var.consensus
  orderer_node_num      = var.orderer_node_num
  cce_cluster_id        = var.cce_cluster_id
  enterprise_project_id = var.enterprise_project_id
  password              = var.instance_password
  volume_type           = var.volume_type
  org_disk_size         = var.org_disk_size

  dynamic "block_info" {
    for_each = var.block_info

    content {
      generation_interval  = block_info.value.generation_interval
      transaction_quantity = block_info.value.transaction_quantity
      block_size           = block_info.value.block_size
    }
  }

  dynamic "sfs_turbo" {
    for_each = var.sfs_turbo

    content {
      share_type        = sfs_turbo.value.share_type
      type              = sfs_turbo.value.type
      availability_zone = sfs_turbo.value.availability_zone
      flavor            = sfs_turbo.value.flavor
    }
  }

  dynamic "peer_orgs" {
    for_each = var.peer_orgs

    content {
      org_name = peer_orgs.value.org_name
      count    = peer_orgs.value.count
    }
  }

  dynamic "channels" {
    for_each = var.channels

    content {
      name      = channels.value.name
      org_names = channels.value.org_names
    }
  }
}
