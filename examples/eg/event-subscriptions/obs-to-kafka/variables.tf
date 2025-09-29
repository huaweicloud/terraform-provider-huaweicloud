# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPC is located"
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
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "172.16.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = "172.16.10.0/24"
}

variable "subnet_gateway" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = "172.16.10.1"
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "bucket_name" {
  description = "The name of the OBS bucket"
  type        = string
}

variable "bucket_acl" {
  description = "The ACL policy for a bucket"
  type        = string
  default     = "private"
}

variable "availability_zones" {
  description = "The availability zones to which the Kafka instance belongs"
  type        = list(string)
  default     = []
}

variable "instance_flavor_id" {
  description = "The flavor ID of the Kafka instance"
  type        = string
  default     = "kafka.2u4g.cluster.small"
}

variable "instance_flavor_type" {
  description = "The flavor type of the Kafka instance"
  type        = string
  default     = "cluster"
}

variable "instance_storage_spec_code" {
  description = "The storage specification code of the Kafka instance"
  type        = string
  default     = "dms.physical.storage.high.v2"
}

variable "instance_name" {
  description = "The name of the Kafka instance"
  type        = string
}

variable "instance_engine_version" {
  description = "The engine version of the Kafka instance"
  type        = string
  default     = "3.x"
}

variable "instance_storage_space" {
  description = "The storage space of the Kafka instance"
  type        = number
  default     = 300
}

variable "instance_broker_num" {
  description = "The number of brokers of the Kafka instance"
  type        = number
  default     = 3
}

variable "instance_ssl_enable" {
  description = "The SSL enable of the Kafka instance"
  type        = bool
  default     = false
}

variable "instance_description" {
  description = "The description of the Kafka instance"
  type        = string
  default     = ""
}

variable "instance_security_protocol" {
  description = "The protocol to use after SASL is enabled"
  type        = string
  default     = "SASL_SSL"
}

variable "charging_mode" {
  description = "The charging mode of the Kafka instance"
  type        = string
  default     = "postPaid"
}

variable "topic_name" {
  description = "The name of the topic"
  type        = string
}

variable "topic_partitions" {
  description = "The number of the topic partition"
  type        = number
  default     = 3
}

variable "connection_name" {
  description = "The name of the connection"
  type        = string
}

variable "connection_acks" {
  description = "The number of confirmation signals the prouder needs to receive to consider the message sent successfully"
  type        = string
  default     = "1"
}

variable "subscription_source_values" {
  description = "The event types to be subscribed from OBS service"
  type        = list(string)
  default     = [
    "OBS:CloudTrace:ApiCall",
    "OBS:CloudTrace:ObsSDK",
    "OBS:CloudTrace:ConsoleAction",
    "OBS:CloudTrace:SystemAction",
    "OBS:CloudTrace:Others"
  ]
}

variable "object_extension_name" {
  description = "The extension name of the OBS object to be uploaded"
  type        = string
  default     = ".txt"
  nullable    = false
}

variable "object_name" {
  description = "The name of the OBS object to be uploaded"
  type        = string
}

variable "object_upload_content" {
  description = "The content of the OBS object to be uploaded"
  type        = string
}
