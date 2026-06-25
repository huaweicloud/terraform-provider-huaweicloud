# Variable definitions for authentication
variable "region_name" {
  type        = string
  description = "The region where the ModelArts service is located"
}

variable "access_key" {
  type        = string
  sensitive   = true
  description = "The access key of the IAM user"
}

variable "secret_key" {
  type        = string
  sensitive   = true
  description = "The secret key of the IAM user"
}

# Variable definitions for resources/data sources
variable "vpc_name" {
  type        = string
  description = "The VPC name"
}

variable "vpc_cidr" {
  type        = string
  default     = "192.168.0.0/16"
  description = "The CIDR block of the VPC"
}

variable "enterprise_project_id" {
  type        = string
  default     = null
  description = "The enterprise project ID"
}

variable "subnet_name" {
  type        = string
  description = "The subnet name"
}

variable "subnet_cidr" {
  type        = string
  default     = ""
  nullable    = false
  description = "The CIDR block of the subnet"
}

variable "subnet_gateway_ip" {
  type        = string
  default     = ""
  nullable    = false
  description = "The gateway IP of the subnet"
}

variable "security_group_name" {
  type        = string
  description = "The name of the security group"
}

variable "turbo_name" {
  type        = string
  description = "The name of the SFS Turbo"
}

variable "turbo_size" {
  type        = number
  default     = 1228
  description = "The size of the SFS Turbo"
}

variable "turbo_share_proto" {
  type        = string
  default     = "NFS"
  description = "The share protocol of the SFS Turbo"
}

variable "turbo_share_type" {
  type        = string
  default     = "HPC"
  description = "The share type of the SFS Turbo"
}

variable "turbo_hpc_bandwidth" {
  type        = string
  default     = "40M"
  description = "The HPC bandwidth of the SFS Turbo"
}

variable "network_name" {
  type        = string
  description = "The name of the network"
}

variable "network_cidr" {
  type        = string
  default     = "10.168.0.0/16"
  description = "The CIDR block of the network"
}

variable "workspace_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The name of the workspace to create. Cannot be configured together with workspace_id."

  validation {
    condition     = var.workspace_name == "" || var.workspace_id == ""
    error_message = "workspace_name and workspace_id cannot be configured at the same time."
  }
}

variable "resource_pool_flavor_id" {
  type        = string
  default     = ""
  nullable    = false
  description = "The flavor ID of the resource pool"
}

variable "resource_pool_name" {
  type        = string
  description = "The name of the resource pool"
}

variable "resource_pool_scope" {
  type        = list(string)
  default     = ["Notebook", "Train", "Infer"]
  description = "The scope of the resource pool"
}

variable "workspace_id" {
  type        = string
  default     = ""
  description = "The existing workspace ID. Cannot be configured together with workspace_name."
}

variable "notebook_flavor_id" {
  type        = string
  default     = ""
  nullable    = false
  description = "The flavor ID of the notebook. If empty, the first available dedicated flavor is used"
}

variable "notebook_flavor_category" {
  type        = string
  default     = "CPU"
  description = "The processor type of the notebook flavor"
}

variable "notebook_image_id" {
  type        = string
  default     = ""
  nullable    = false
  description = "The image ID of the notebook."
}

variable "notebook_image_type" {
  type        = string
  default     = "BUILD_IN"
  description = "The type of the notebook image"
}

variable "notebook_key_pair_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The key pair name for remote SSH access"

  validation {
    condition     = length(var.allowed_access_ips) == 0 || var.notebook_key_pair_name != "" || var.keypair_name != ""
    error_message = "When allowed_access_ips is configured, either notebook_key_pair_name or keypair_name must be specified."
  }
}

variable "allowed_access_ips" {
  type        = list(string)
  default     = []
  nullable    = false
  description = "The public IP addresses that are allowed for remote SSH access"
}

variable "keypair_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The name of the KPS key pair"
}

variable "notebook_name" {
  type        = string
  description = "The name of the notebook"
}

variable "notebook_description" {
  type        = string
  default     = null
  description = "The description of the notebook"
}

variable "notebook_tags" {
  type        = map(string)
  default     = {}
  description = "The tags of the notebook"
}

variable "notebook_mount_storage_path" {
  type        = string
  default     = ""
  nullable    = false
  description = "The OBS path of Parallel File System (PFS) or its folders to mount"
}

variable "notebook_mount_storage_local_directory" {
  type        = string
  default     = ""
  nullable    = false
  description = "The local mount directory for the OBS storage"
}
