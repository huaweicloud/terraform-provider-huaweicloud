# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ECS instance is located"
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
variable "availability_zone" {
  description = "The availability zone for the ECS instance"
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
  description = "The number of CPU cores for the ECS instance"
  type        = number
  default     = 2
}

variable "instance_memory_size" {
  description = "The memory size in GB for the ECS instance"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The operating system of the ECS instance image"
  type        = string
  default     = "Ubuntu"
}

variable "keypair_name" {
  description = "The name of the SSH keypair"
  type        = string
}

variable "private_key_path" {
  description = "The path to the private key file"
  type        = string
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

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
  default     = ""
}

variable "associate_eip_address" {
  description = "The EIP address to associate with the ECS instance"
  type        = string
  default     = ""
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
  default     = ""
}

variable "bandwidth_size" {
  description = "The size of the bandwidth"
  type        = number
  default     = 5
}

variable "bandwidth_share_type" {
  description = "The share type of the bandwidth"
  type        = string
  default     = "PER"
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the bandwidth"
  type        = string
  default     = "traffic"
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "instance_user_data" {
  description = "The user data script for the ECS instance"
  type        = string
  default     = ""
}

variable "instance_system_disk_type" {
  description = "The system disk type for the ECS instance"
  type        = string
  default     = "SSD"
}

variable "instance_system_disk_size" {
  description = "The system disk size in GB for the ECS instance"
  type        = number
  default     = 40
}

variable "security_group_ids" {
  description = "The list of security group IDs for the ECS instance"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "instance_remote_exec_inline" {
  description = "The inline script for remote execution"
  type        = list(string)
  nullable    = false
}
