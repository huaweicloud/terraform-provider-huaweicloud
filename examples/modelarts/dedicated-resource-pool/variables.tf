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
variable "turbo_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The name of the SFS Turbo to create"

  validation {
    condition     = var.network_name != "" || var.turbo_name == ""
    error_message = "turbo_name can only be specified when network_name is specified."
  }
}

variable "vpc_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The VPC name. Can only be specified when turbo_name is specified."

  validation {
    condition     = (var.turbo_name != "") == (var.vpc_name != "")
    error_message = "vpc_name can only be specified when turbo_name is specified."
  }
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
  default     = ""
  nullable    = false
  description = "The subnet name. Can only be specified when turbo_name is specified."

  validation {
    condition     = (var.turbo_name != "") == (var.subnet_name != "")
    error_message = "subnet_name can only be specified when turbo_name is specified."
  }
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
  default     = ""
  nullable    = false
  description = "The name of the security group"

  validation {
    condition     = (var.turbo_name != "") == (var.security_group_name != "")
    error_message = "security_group_name can only be specified when turbo_name is specified."
  }
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

variable "network_sfs_turbos" {
  type = list(object({
    id   = string
    name = string
  }))

  default     = []
  nullable    = false
  description = "The SFS Turbo connections for the ModelArts network"

  validation {
    condition     = length(var.network_sfs_turbos) == 0 || (var.network_name != "" && var.turbo_name == "")
    error_message = "network_sfs_turbos can only be specified when network_name is specified and turbo_name is not specified."
  }
}

variable "network_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The name of the ModelArts network to create"

  validation {
    condition     = (var.network_id != "") != (var.network_name != "")
    error_message = "Exactly one of network_id and network_name must be specified."
  }
}

variable "network_cidr" {
  type        = string
  default     = "10.168.0.0/16"
  description = "The CIDR block of the ModelArts network"
}

variable "workspace_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The name of the workspace to create"

  validation {
    condition     = var.workspace_name == "" || var.workspace_id == ""
    error_message = "workspace_name and workspace_id cannot be configured at the same time."
  }
}

variable "resource_pool_resources" {
  type = list(object({
    flavor_id     = optional(string, "")
    count         = number
    max_count     = optional(number)
    extend_params = optional(string)

    root_volume = optional(object({
      volume_type = string
      size        = string
    }))

    data_volumes = optional(list(object({
      volume_type   = string
      size          = string
      count         = optional(number)
      extend_params = optional(string)
    })), [])

    volume_group_configs = optional(list(object({
      volume_group     = string
      docker_thin_pool = optional(number)
      types            = optional(list(string))

      lvm_config = optional(object({
        lv_type = string
        path    = optional(string)
      }))
    })), [])

    os = optional(object({
      name       = optional(string)
      image_id   = optional(string)
      image_type = optional(string)
    }))

    driver = optional(object({
      version = string
    }))

    creating_step = optional(object({
      step = number
      type = string
    }))
  }))

  description = "The list of resource specifications in the resource pool"
}

variable "resource_pool_name" {
  type        = string
  description = "The name of the dedicated resource pool"
}

variable "resource_pool_description" {
  type        = string
  default     = null
  description = "The description of the dedicated resource pool"
}

variable "resource_pool_scope" {
  type        = list(string)
  default     = ["Train", "Infer", "Notebook"]
  description = "The scope of the dedicated resource pool"
}

variable "network_id" {
  type        = string
  default     = ""
  nullable    = false
  description = "The existing ModelArts network ID"
}

variable "workspace_id" {
  type        = string
  default     = ""
  description = "The existing workspace ID"
}

variable "resource_pool_metadata_annotations" {
  type        = string
  default     = null
  description = "The annotations of the resource pool, in JSON format"
}
