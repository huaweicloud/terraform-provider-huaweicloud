# Variable definitions for authentication
variable "region_name" {
  description = "The region where the Kafka instance is located"
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

# Variable definitions for resources/data sources
variable "instance_configurations" {
  description = "The list of configurations for multiple Kafka instances"

  type = list(object({
    name               = string
    availability_zones = optional(list(string), [])
    engine_version     = optional(string, "3.x")
    flavor_id          = optional(string, "")
    flavor_type        = optional(string, "cluster")
    storage_spec_code  = optional(string, "dms.physical.storage.ultra.v2")
    storage_space      = optional(number, 600)
    broker_num         = optional(number, 3)
    access_user        = optional(string, "")
    password           = optional(string, "")
    enabled_mechanisms = optional(list(string), null)

    port_protocol = optional(object({
      private_plain_enable          = optional(bool, true)
      private_sasl_ssl_enable       = optional(bool, null)
      private_sasl_plaintext_enable = optional(bool, null)
    }), {})
  }))

  nullable = false
  default  = []

  validation {
    condition     = length(var.instance_configurations) >= 2
    error_message = "At least two instances are required"
  }
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "task_topics" {
  description = "The topics of the Smart Connect task"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "topic_name" {
  description = "The name of the Kafka topic"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.topic_name != "" || length(var.task_topics) > 0
    error_message = "topic_name is required when task_topics is not provided"
  }
}

variable "topic_partitions" {
  description = "The number of partitions of the topic"
  type        = number
  default     = 10
}

variable "topic_replicas" {
  description = "The number of replicas of the topic"
  type        = number
  default     = 3
}

variable "topic_aging_time" {
  description = "The aging time of the topic"
  type        = number
  default     = 72
}

variable "topic_sync_replication" {
  description = "The sync replication of the topic"
  type        = bool
  default     = false
}

variable "topic_sync_flushing" {
  description = "The sync flushing of the topic"
  type        = bool
  default     = false
}

variable "topic_description" {
  description = "The description of the topic"
  type        = string
  default     = null
}

variable "topic_configs" {
  description = "The configs of the topic"

  type = list(object({
    name  = string
    value = string
  }))

  default  = []
  nullable = false
}

variable "smart_connect_storage_spec_code" {
  description = "The storage specification code of the Smart Connect"
  type        = string
  default     = null
}

variable "smart_connect_bandwidth" {
  description = "The bandwidth of the Smart Connect"
  type        = string
  default     = null
}

variable "smart_connect_node_count" {
  description = "The number of nodes of the Smart Connect"
  type        = number
  default     = 2
}

variable "task_name" {
  description = "The name of the Smart Connect task"
  type        = string
}

variable "task_start_later" {
  description = "The start later of the Smart Connect task"
  type        = bool
  default     = false
}

variable "task_direction" {
  description = "The direction of the Smart Connect task"
  type        = string
  default     = "two-way"
}

variable "task_replication_factor" {
  description = "The replication factor of the Smart Connect task"
  type        = number
  default     = 3
}

variable "task_task_num" {
  description = "The number of tasks of the Smart Connect task"
  type        = number
  default     = 2
}

variable "task_provenance_header_enabled" {
  description = "The provenance header enabled of the Smart Connect task"
  type        = bool
  default     = false
}

variable "task_sync_consumer_offsets_enabled" {
  description = "The sync consumer offsets enabled of the Smart Connect task"
  type        = bool
  default     = false
}

variable "task_rename_topic_enabled" {
  description = "The rename topic enabled of the Smart Connect task"
  type        = bool
  default     = true
}

variable "task_consumer_strategy" {
  description = "The consumer strategy of the Smart Connect task"
  type        = string
  default     = "latest"
}

variable "task_compression_type" {
  description = "The compression type of the Smart Connect task"
  type        = string
  default     = "none"
}

variable "task_topics_mapping" {
  description = "The topics mapping of the Smart Connect task"
  type        = list(string)
  default     = []
}
