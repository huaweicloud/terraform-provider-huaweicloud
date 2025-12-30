# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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
  description = "The availability zone where the resources will be created"
  type        = string
  default     = ""
  nullable    = false
}

variable "deh_instance_host_type" {
  description = "The host type of the dedicated host"
  type        = string
  default     = ""
  nullable    = false
}

variable "deh_instance_name" {
  description = "The name of the dedicated host instance"
  type        = string
}

variable "deh_instance_auto_placement" {
  description = "Whether to enable auto placement for the dedicated host"
  type        = string
  default     = "on"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the dedicated host"
  type        = string
  default     = null
}

variable "deh_instance_charging_mode" {
  description = "The charging mode of the dedicated host"
  type        = string
  default     = "prePaid"
}

variable "deh_instance_period_unit" {
  description = "The unit of the billing period of the dedicated host"
  type        = string
  default     = "month"
}

variable "deh_instance_period" {
  description = "The billing period of the dedicated host"
  type        = string
  default     = "1"
}

variable "deh_instance_auto_renew" {
  description = "Whether to enable auto renew for the dedicated host"
  type        = string
  default     = "false"
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
  description = "The name of the VPC subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the VPC subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the VPC subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "ecs_instance_image_id" {
  description = "The ID of the ECS instance image"
  type        = string
  default     = ""
  nullable    = false
}

variable "ecs_instance_flavor_id" {
  description = "The ID of the ECS instance flavor"
  type        = string
  default     = ""
  nullable    = false
}

variable "ecs_instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "ecs_instance_image_os" {
  description = "The OS of the ECS instance image"
  type        = string
  default     = "Ubuntu"
}

variable "ecs_instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "ecs_instance_admin_pass" {
  description = "The password of the ECS instance administrator"
  type        = string
  sensitive   = true
}
