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

# Variable definitions for resources and data sources
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

variable "topic_name" {
  type        = string
  default     = ""
  nullable    = false
  description = "The name of the SMN topic to create for notifications. Cannot be configured together with training_job_notification_topic_urn."

  validation {
    condition     = var.topic_name == "" || var.training_job_notification_topic_urn == ""
    error_message = "topic_name and training_job_notification_topic_urn cannot be configured at the same time."
  }
}

variable "training_job_name" {
  type        = string
  description = "The name of the training job"
}

variable "training_job_annotations" {
  type        = map(string)
  default     = {}
  description = "The annotations of the training job"
}

variable "training_job_description" {
  type        = string
  default     = null
  description = "The description of the training job"
}

variable "training_job_code_dir" {
  type        = string
  description = "The OBS code directory of the training job"
}

variable "training_job_command" {
  type        = string
  description = "The container startup command for the training job"
}

variable "training_job_engine" {
  type = object({
    image_url = optional(string)
    id        = optional(string)
    version   = optional(string)
    name      = optional(string)
  })

  description = "The engine configuration of the training job"
}

variable "training_job_inputs" {
  type = list(object({
    local_dir = string

    dataset = object({
      id           = string
      name         = optional(string)
      version_id   = optional(string)
      service_type = optional(string)
    })
  }))

  default     = []
  nullable    = false
  description = "The inputs of the training job"
}

variable "training_job_environments" {
  type        = map(string)
  default     = {}
  nullable    = false
  description = "The environment variables of the training job"
}

variable "resource_node_count" {
  type        = number
  default     = 1
  description = "The number of resource replicas used by the training job"
}

variable "training_job_volumes" {
  type = list(object({
    nfs = optional(object({
      nfs_server_path = optional(string)
      local_path      = optional(string)
      read_only       = optional(bool)
    }))

    pfs = optional(object({
      pfs_path   = optional(string)
      local_path = optional(string)
    }))

    obs = optional(object({
      obs_path   = optional(string)
      local_path = optional(string)
    }))
  }))

  default     = []
  nullable    = false
  description = "The volume mount configuration of the training job"
}

variable "training_job_log_export_path_obs_url" {
  type        = string
  default     = ""
  nullable    = false
  description = "The log export path of the training job"
}

variable "training_job_log_export_config_version" {
  type        = string
  default     = ""
  nullable    = false
  description = "The log export config version of the training job"
}

variable "training_job_auto_stop_duration" {
  type        = number
  default     = 0
  description = "The auto stop duration of the training job"
}

variable "training_job_notification_topic_urn" {
  type        = string
  default     = ""
  nullable    = false
  description = "The existing SMN topic URN for training job notifications. Cannot be configured together with topic_name."
}

variable "training_job_notification_events" {
  type        = list(string)
  default     = []
  description = "The notification events of the training job"
}

variable "custom_metrics" {
  type = object({
    exec = optional(object({
      command = list(string)
    }))

    http_get = optional(object({
      path = string
      port = number
    }))
  })

  default     = null
  description = "The custom metrics of the training job"
}

variable "training_job_asset_model" {
  type = object({
    name    = string
    version = string
    type    = string
    code    = optional(string)
    desc    = optional(string)
    series  = optional(string)
  })

  default     = null
  description = "The asset model of the training job"
}

variable "training_job_output_model" {
  type = object({
    obs_path   = string
    local_path = optional(string)
  })

  default     = null
  description = "The output model of the training job"
}

variable "training_job_tags" {
  type        = map(string)
  default     = {}
  description = "The tags of the training job"
}
