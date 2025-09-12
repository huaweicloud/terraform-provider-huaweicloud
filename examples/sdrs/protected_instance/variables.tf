# Variable definitions for authentication
variable "region_name" {
  description = "The region where the protection instance is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for resources/data sources
variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
}

variable "availability_zone" {
  description = "The availability zone to which the ECS instance and network belong"
  type        = string
  default     = ""
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The number of the ECS instance flavor CPU cores"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The number of the ECS instance flavor memories"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
}

variable "instance_image_os_type" {
  description = "The OS type of the ECS instance flavor"
  type        = string
  default     = "Ubuntu"
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance flavor"
  type        = string
  default     = "public"
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
  default     = "192.168.0.0/24"
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "ecs_instance_name" {
  type        = string
  description = "The name of the ECS instance"
}

variable "sdrs_domain_name" {
  description = "The name of an available SDRS domain"
  type        = string
  default     = null
}

variable "protection_group_name" {
  description = "The name of the protection group"
  type        = string
}

variable "source_availability_zone" {
  description = "The production site AZ of the protection group"
  type        = string
  default     = ""

  validation {
    condition = (
      (var.source_availability_zone == "" && var.target_availability_zone == "") ||
      (var.source_availability_zone != "" && var.target_availability_zone != "")
    )
    error_message = "Both `source_availability_zone` and `target_availability_zone` must be set, or both must be empty"
  }
}

variable "target_availability_zone" {
  description = "The disaster recovery site AZ of the protection group"
  type        = string
  default     = ""
}

variable "protected_instance_name" {
  description = "The name of the protected instance"
  type        = string
}

variable "cluster_id" {
  description = "The DSS storage pool ID"
  type        = string
  default     = null
}

variable "primary_ip_address" {
  description = "The IP address of the primary NIC on the DR site server"
  type        = string
  default     = "192.168.0.15"
}

variable "delete_target_server" {
  description = "Whether to delete the DR site server"
  type        = bool
  default     = null
}

variable "delete_target_eip" {
  description = "Whether to delete the EIP of the DR site server"
  type        = bool
  default     = null
}

variable "protected_instance_description" {
  description = "The description of the protected instance"
  type        = string
  default     = null
}

variable "protected_instance_tags" {
  description = "The key/value pairs to associate with the protected instance"
  type        = map(string)
  default     = null
}
