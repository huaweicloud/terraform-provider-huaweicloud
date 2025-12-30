# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
  type        = string
}

# Variable definitions for sharer account
variable "access_key" {
  description = "The access key of the IAM user in sharer account"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user in sharer account"
  type        = string
  sensitive   = true
}

# Variable definitions for accepter account
variable "accepter_access_key" {
  description = "The access key of the IAM user in accepter account"
  type        = string
  sensitive   = true
}

variable "accepter_secret_key" {
  description = "The secret key of the IAM user in accepter account"
  type        = string
  sensitive   = true
}

variable "instance_flavor_id" {
  description = "The ID of the ECS instance flavor"
  type        = string
  default     = ""
  nullable    = true
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the ECS instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the ECS instance flavor (GB)"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The ID of the ECS instance image"
  type        = string
  default     = ""
  nullable    = true
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The OS of the ECS instance image"
  type        = string
  default     = "Ubuntu"
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

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = null
}

variable "subnet_name" {
  description = "The name of the VPC subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the VPC subnet"
  type        = string
  default     = ""
  nullable    = true
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the VPC subnet"
  type        = string
  default     = ""
  nullable    = true
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  description = "The name of the ECS instance to be created"
  type        = string
}

variable "administrator_password" {
  description = "The password of the administrator for the ECS instance"
  type        = string
  sensitive   = true
}

variable "data_volume_name" {
  description = "The name of the data volume to be created and attached to ECS instance"
  type        = string
}

variable "data_volume_type" {
  description = "The type of the data volume"
  type        = string
  default     = "SAS"
}

variable "data_volume_size" {
  description = "The size of the data volume in GB"
  type        = number
  default     = 10
}

# Variable definitions for data image
variable "data_image_name" {
  description = "The name of the data disk image to be created"
  type        = string
}

variable "data_image_description" {
  description = "The description of the data disk image"
  type        = string
  default     = ""
}

variable "accepter_vpc_name" {
  description = "The name of the VPC in accepter account"
  type        = string
}

variable "accepter_vpc_cidr" {
  description = "The CIDR block of the VPC in accepter account"
  type        = string
  default     = "192.168.0.0/16"
}

variable "accepter_subnet_name" {
  description = "The name of the VPC subnet in accepter account"
  type        = string
}

variable "accepter_security_group_name" {
  description = "The name of the security group in accepter account"
  type        = string
}

variable "accepter_instance_flavor_id" {
  description = "The ID of the ECS instance flavor in accepter account"
  type        = string
  default     = ""
  nullable    = true
}

variable "accepter_instance_image_id" {
  description = "The ID of the ECS instance image in accepter account"
  type        = string
  default     = ""
  nullable    = true
}

variable "accepter_instance_name" {
  description = "The name of the ECS instance to be created in accepter account (optional, for creating new instance)"
  type        = string
}

variable "accepter_data_volume_name" {
  description = "The name of the data volume to be created from shared image in accepter account"
  type        = string
}

variable "accepter_data_volume_type" {
  description = "The type of the data volume in accepter account"
  type        = string
  default     = "SAS"
}

variable "accepter_data_volume_size" {
  description = "The size of the data volume in accepter account in GB"
  type        = number
  default     = 20
}
