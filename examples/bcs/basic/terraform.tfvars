# BCS instance configuration
instance_name         = "wzp-test3"
edition               = 4
fabric_version        = "4.0.35"
consensus             = "etcdraft"
orderer_node_num      = 3
cce_cluster_id        = "1081aa42-7b69-11f1-946d-0255ac100260"
enterprise_project_id = "571f1c69-6d32-43d9-8ba5-eff73fc3eebe"
instance_password     = "SLKop123S@%!dsa1"
volume_type           = "efs"
org_disk_size         = 3686

block_info = [
  {
    generation_interval  = 2
    transaction_quantity = 500
    block_size           = 2
  }
]

# SFS Turbo configuration (required for volume_type = "efs" and edition = 4)
sfs_turbo = [
  {
    share_type        = "STANDARD"
    type              = "efs-ha"
    availability_zone = "cn-north-4a"
    flavor            = "sfs.turbo.20MBps"
  }
]

# Peer organizations configuration
peer_orgs = [
  {
    org_name = "organization"
    count    = 2
  }
]

# Channels configuration
channels = [
  {
    name      = "channel"
    org_names = ["organization"]
  }
]
