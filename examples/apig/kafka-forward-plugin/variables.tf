# Variable definitions for authentication
variable "region_name" {
  description = "The name of the region"
  type        = string
}

variable "access_key" {
  description = "The access key for the Huawei Cloud"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key for the Huawei Cloud"
  type        = string
  sensitive   = true
}

# Variable definitions for resources/data sources
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
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "availability_zones" {
  description = "The availability zones to which the instance belongs"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "instance_name" {
  description = "The name of the APIG instance"
  type        = string
}

variable "instance_edition" {
  description = "The edition of the APIG instance"
  type        = string
  default     = "BASIC"
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

variable "availability_zones_count" {
  description = "The number of availability zones to which the instance belongs"
  type        = number
  default     = 1
}

variable "kafka_instance_flavor_id" {
  description = "The flavor ID of the DMS Kafka instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "kafka_instance_flavor_type" {
  description = "The flavor type of the DMS Kafka instance"
  type        = string
  default     = "cluster"
}

variable "kafka_instance_storage_spec_code" {
  description = "The storage spec code of the DMS Kafka instance"
  type        = string
  default     = "dms.physical.storage.high.v2"
}

variable "kafka_instance_name" {
  description = "The name of the DMS Kafka instance"
  type        = string
}

variable "kafka_instance_description" {
  description = "The description of the DMS Kafka instance"
  type        = string
  default     = ""
}

variable "kafka_instance_ssl_enable" {
  description = "Whether to enable SSL for the DMS Kafka instance"
  type        = bool
  default     = false
}

variable "kafka_instance_engine_version" {
  description = "The engine version of the DMS Kafka instance"
  type        = string
}

variable "kafka_instance_storage_space" {
  description = "The storage space of the DMS Kafka instance in GB"
  type        = number
}

variable "kafka_instance_broker_num" {
  description = "The number of brokers for the DMS Kafka instance"
  type        = number
}

variable "kafka_charging_mode" {
  description = "The charging mode of the DMS Kafka instance. Options: prePaid, postPaid"
  type        = string
  default     = "prePaid"
}

variable "kafka_period_unit" {
  description = "The period unit of the DMS Kafka instance. Options: month, year"
  type        = string
  default     = "month"
}

variable "kafka_period" {
  description = "The period of the DMS Kafka instance"
  type        = number
  default     = 1
}

variable "kafka_auto_new" {
  description = "Whether to enable auto renewal for the DMS Kafka instance"
  type        = string
  default     = "false"
}

variable "kafka_instance_user_name" {
  description = "The access user name for the DMS Kafka instance"
  type        = string
  sensitive   = true
}

variable "kafka_instance_user_password" {
  description = "The access user password for the DMS Kafka instance"
  type        = string
  sensitive   = true
}

variable "kafka_topic_name" {
  description = "The name of the Kafka topic to receive messages"
  type        = string
}

variable "kafka_topic_partitions" {
  description = "The number of partitions for the Kafka topic"
  type        = number
  default     = 1
}

variable "plugin_name" {
  description = "The name of the Kafka forward plugin"
  type        = string
}

variable "plugin_description" {
  description = "The description of the Kafka forward plugin"
  type        = string
  default     = null
}

variable "kafka_security_protocol" {
  description = "The security protocol for Kafka connection. Options: PLAINTEXT, SASL_PLAINTEXT, SASL_SSL, SSL"
  type        = string
  default     = "PLAINTEXT"
  nullable    = false

  validation {
    condition     = contains(["PLAINTEXT", "SASL_PLAINTEXT", "SASL_SSL", "SSL"], var.kafka_security_protocol)
    error_message = "kafka_security_protocol must be one of: PLAINTEXT, SASL_PLAINTEXT, SASL_SSL, SSL."
  }
}

variable "kafka_message_key" {
  description = "The message key extraction strategy. Can be a static value or a variable expression like $context.requestId"
  type        = string
  default     = ""
}

variable "kafka_max_retry_count" {
  description = "The maximum number of retry attempts for failed message sends"
  type        = number
  default     = 3
}

variable "kafka_retry_backoff" {
  description = "The backoff time in seconds between retries"
  type        = number
  default     = 10
}

variable "kafka_sasl_mechanisms" {
  description = "The SASL mechanism for authentication. Options: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512"
  type        = string
  default     = "PLAIN"

  validation {
    condition     = contains(["PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512"], var.kafka_sasl_mechanisms)
    error_message = "kafka_sasl_mechanisms must be one of: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512."
  }
}

variable "kafka_sasl_username" {
  description = "The SASL username for authentication (leave empty to use kafka_access_user)"
  type        = string
  default     = ""
  sensitive   = true
  nullable    = false
}

variable "kafka_access_user" {
  description = "The access user for Kafka authentication (used when kafka_sasl_username is empty and security_protocol is not PLAINTEXT)"
  type        = string
  default     = ""
  sensitive   = true
  nullable    = false
}

variable "kafka_sasl_password" {
  description = "The SASL password for authentication (leave empty to use kafka_password)"
  type        = string
  default     = ""
  sensitive   = true
  nullable    = false
}

variable "kafka_password" {
  description = "The password for Kafka authentication (used when kafka_sasl_password is empty and security_protocol is not PLAINTEXT)"
  type        = string
  default     = ""
  sensitive   = true
  nullable    = false
}

variable "kafka_ssl_ca_content" {
  description = "The SSL CA certificate content for SSL/TLS encrypted connections"
  type        = string
  default     = ""
  sensitive   = true
  nullable    = false
}
