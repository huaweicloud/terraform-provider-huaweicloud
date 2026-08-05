# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CSS cluster is located"
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
variable "availability_zone" {
  description = "The availability zone to which the CSS cluster belongs"
  type        = string
  default     = ""
  nullable    = false
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

variable "cluster_flavor" {
  description = "The flavor of the CSS cluster"
  type        = string
  default     = ""
  nullable    = false
}

variable "cluster_name" {
  description = "The name of the CSS cluster"
  type        = string
}

variable "cluster_engine_version" {
  description = "The engine version of the CSS cluster"
  type        = string
  default     = "7.10.2"
}

variable "cluster_instance_number" {
  description = "The number of instances of the CSS cluster"
  type        = number
  default     = 3
}

variable "cluster_volume_type" {
  description = "The volume type of the CSS cluster"
  type        = string
  default     = "ULTRAHIGH"
}

variable "cluster_volume_size" {
  description = "The volume size of the CSS cluster"
  type        = number
  default     = 40
}

variable "ai_ops_check_type" {
  description = "The check type of the AI Ops setting"
  type        = string
  default     = "full_detection"
}

variable "ai_ops_period" {
  description = "The period of the AI Ops setting"
  type        = string
  default     = "12:00 GMT+08:00"
}

variable "ai_ops_check_items" {
  description = "The check items of the AI Ops setting"
  type        = list(string)
  default     = null
}
