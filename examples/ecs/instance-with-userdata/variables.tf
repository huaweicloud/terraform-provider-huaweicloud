# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ECS instance is located"
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
  description = "The availability zone to which the ECS instance belongs"
  type        = string
  default     = ""
}

variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
}

variable "instance_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_cpu_core_count" {
  description = "The number of the vCPUs in the ECS instance flavor"
  type        = number
  default     = 2
}

variable "instance_memory_size" {
  description = "The memory size(GB) in the ECS instance flavor"
  type        = number
  default     = 4
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
}

variable "subnet_gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "security_group_ids" {
  description = "The list of security group IDs for the ECS instance"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "security_group_names" {
  description = "The name of the security groups to which the ECS instance belongs"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
}

variable "instance_image_visibility" {
  description = "The visibility of the image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The operating system of the image"
  type        = string
  default     = "Ubuntu"
}

variable "keypair_name" {
  description = "The keypair name for the ECS instance"
  type        = string
}

variable "keypair_public_key" {
  description = "The public key for the keypair"
  type        = string
  default     = null
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "instance_user_data" {
  description = "The user data script for the ECS instance"
  type        = string
}
