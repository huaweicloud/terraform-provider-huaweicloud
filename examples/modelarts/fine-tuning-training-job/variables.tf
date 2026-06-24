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

variable "training_job_train_type" {
  type        = string
  default     = "SFT"
  description = "The training type of the fine-tuning job"
}

variable "training_job_name" {
  type        = string
  description = "The name of the training job"
}

variable "workspace_id" {
  type        = string
  default     = ""
  description = "The existing workspace ID. Cannot be configured together with workspace_name."
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

variable "training_job_inputs" {
  type = list(object({
    dataset = object({
      id                 = string
      name               = string
      version_id         = optional(string)
      service_type       = optional(string)
      dataset_proportion = optional(number)
    })

    local_dir = optional(string)
  }))

  default     = []
  nullable    = false
  description = "The inputs of the training job"
}

variable "training_job_environments" {
  type        = map(string)
  default     = null
  description = "The environment variables of the training job"
}

variable "resource_flavor_id" {
  type        = string
  default     = ""
  nullable    = false
  description = "The flavor ID of the public resource pool"
}

variable "resource_node_count" {
  type        = number
  default     = 1
  description = "The number of resource replicas used by the training job"
}

variable "training_job_asset_id" {
  type        = string
  description = "The asset ID of the training job"
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

  description = "The asset model configuration of the training job"
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

  validation {
    condition     = (var.training_job_notification_topic_urn == "" && var.topic_name == "") || length(var.training_job_notification_events) > 0
    error_message = "training_job_notification_events is required when training_job_notification_topic_urn or topic_name is configured."
  }
}

variable "training_job_output_model" {
  type = object({
    obs_path   = string
    local_path = optional(string)
  })

  default     = null
  description = "The output model configuration of the training job"
}

variable "training_job_ftjob_config" {
  type = object({
    envs = list(object({
      env_name    = string
      env_type    = string
      value       = string
      label       = optional(string)
      des         = optional(string)
      modifiable  = optional(bool)
      displayable = optional(bool)
    }))

    checkpoint_config = optional(object({
      checkpoint_id        = optional(string)
      save_checkpoints_max = optional(number)
      skipped_steps        = optional(number)
      restore_training     = optional(number)
    }))
  })

  description = "The fine-tuning training job configuration"
}

variable "training_job_tags" {
  type        = map(string)
  default     = null
  description = "The tags of the training job"
}
