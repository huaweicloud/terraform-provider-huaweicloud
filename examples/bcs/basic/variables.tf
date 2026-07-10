# Variable definitions for authentication
variable "region_name" {
  description = "The region where the BCS instance is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for BCS instance
variable "instance_name" {
  description = "The unique name of the BCS instance"
  type        = string
}

variable "edition" {
  description = "The service edition of the BCS instance"
  type        = number
}

variable "fabric_version" {
  description = "The version of fabric for the BCS instance"
  type        = string
}

variable "consensus" {
  description = "The consensus algorithm used by the BCS instance"
  type        = string
}

variable "orderer_node_num" {
  description = "The number of peers in the orderer organization"
  type        = number
  default     = 1
}

variable "cce_cluster_id" {
  description = "The ID of the CCE cluster to attach to the BCS instance"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project that the BCS instance belongs to"
  type        = string
}

variable "instance_password" {
  description = "The resource access and blockchain management password"
  type        = string
  sensitive   = true
}

variable "volume_type" {
  description = "The storage volume type to attach to each organization of the BCS instance"
  type        = string
  default     = "nfs"
}

variable "org_disk_size" {
  description = "The storage capacity of peer organization"
  type        = number
  default     = 100
}

# Variable definitions for block information
variable "block_info" {
  description = "The configuration of block generation"
  type        = list(object({
    generation_interval  = optional(number, 2)
    transaction_quantity = optional(number, 500)
    block_size           = optional(number, 2)
  }))

  default = []
}

# Variable definitions for SFS Turbo
variable "sfs_turbo" {
  description = "The SFS Turbo configuration for BCS instance"
  type        = list(object({
    share_type        = optional(string, "STANDARD")
    type              = optional(string, "efs-ha")
    availability_zone = optional(string, "")
    flavor            = optional(string, "sfs.turbo.20MBps")
  }))

  default = []

  validation {
    condition     = var.volume_type != "efs" || var.edition != 4 || length(var.sfs_turbo) > 0
    error_message = "When using \"volume_type = \"efs\"\" with \"edition = 4\", you must configure \"sfs_turbo\"."
  }
}

# Variable definitions for peer organizations
variable "peer_orgs" {
  description = "The array of one or more peer organizations to attach to the BCS instance"
  type        = list(object({
    org_name = string
    count    = number
  }))

  default = [
    {
      org_name = "organization"
      count    = 2
    }
  ]
}

# Variable definitions for channels
variable "channels" {
  description = "The array of one or more channels to attach to the BCS instance"
  type        = list(object({
    name      = string
    org_names = list(string)
  }))

  default = [
    {
      name      = "channel"
      org_names = ["organization"]
    }
  ]
}
