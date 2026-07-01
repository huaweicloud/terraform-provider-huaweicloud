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
variable "service_group_pool_id" {
  type        = string
  description = "The ID of the dedicated resource pool for the instance group"
}

variable "service_name" {
  type        = string
  description = "The name of the online inference service"
}

variable "service_version" {
  type        = string
  description = "The version of the online inference service"
}

variable "service_type" {
  type        = string
  default     = "REAL_TIME"
  description = "The type of the service"
}

variable "service_deploy_type" {
  type        = string
  default     = null
  description = "The deploy type of the service"
}

variable "service_description" {
  type        = string
  default     = ""
  description = "The description of the online inference service"
}

# Variable definitions for group configuration
variable "service_group_framework" {
  type        = string
  default     = "COMMON"
  description = "The algorithm framework of the instance group"
}

variable "service_group_name" {
  type        = string
  description = "The name of the instance group"
}

variable "service_group_weight" {
  type        = number
  default     = 100
  description = "The weight percentage of the instance group"
}

variable "service_group_count" {
  type        = number
  default     = 1
  description = "The number of service instances in the deployment scenario"
}

variable "service_unit_configs" {
  type = list(object({
    image = object({
      source   = string
      swr_path = string
      id       = optional(string)
    })

    role = optional(string)

    custom_spec = optional(object({
      memory = number
      cpu    = optional(number)
      gpu    = optional(number)
      ascend = optional(number)
    }))

    flavor = optional(string)

    models = optional(list(object({
      source     = string
      mount_path = string
      address    = optional(string)
      source_id  = optional(string)
    })), [])

    codes = optional(list(object({
      source     = string
      mount_path = string
      address    = optional(string)
      source_id  = optional(string)
    })), [])

    count = optional(number)
    cmd   = optional(string)
    envs  = optional(map(string))

    readiness_health = optional(object({
      initial_delay_seconds = number
      timeout_seconds       = number
      period_seconds        = number
      failure_threshold     = number
      check_method          = string
      command               = optional(string)
      url                   = optional(string)
    }))

    startup_health = optional(object({
      initial_delay_seconds = number
      timeout_seconds       = number
      period_seconds        = number
      failure_threshold     = number
      check_method          = string
      command               = optional(string)
      url                   = optional(string)
    }))

    liveness_health = optional(object({
      initial_delay_seconds = number
      timeout_seconds       = number
      period_seconds        = number
      failure_threshold     = number
      check_method          = string
      command               = optional(string)
      url                   = optional(string)
    }))

    port     = optional(number)
    recovery = optional(string)
  }))

  description = "The unit configurations of the instance group"
}

variable "service_runtime_config" {
  type        = string
  description = "The runtime configuration of the service, in JSON format"
}

variable "service_upgrade_config" {
  type        = string
  description = "The upgrade configuration of the service, in JSON format"
}

variable "service_log_configs" {
  type = list(object({
    type          = string
    log_group_id  = optional(string)
    log_stream_id = optional(string)
  }))

  default  = []
  nullable = false

  description = "The log configurations of the service"
}

variable "service_tags" {
  type    = map(string)
  default = {}

  description = "The key/value tags to associate with the service"
}
