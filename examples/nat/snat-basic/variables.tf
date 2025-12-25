# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SNAT rule is located"
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

# Variable definitions for VPC
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
}

# Variable definitions for Subnet
variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

# Variable definitions for NAT Gateway
variable "nat_gateway_name" {
  description = "The name of the NAT gateway"
  type        = string
}

variable "gateway_description" {
  description = "The description of the NAT gateway"
  type        = string
  default     = ""
}

variable "gateway_spec" {
  description = "The specification of the NAT gateway"
  type        = string
  default     = "1"

  validation {
    condition     = contains(["1", "2", "3", "4"], var.gateway_spec)
    error_message = "The gateway_spec must be one of: 1, 2, 3, 4."
  }
}

# Variable definitions for EIP
variable "eip_bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
}

variable "eip_bandwidth_size" {
  description = "The size of the EIP bandwidth in Mbps"
  type        = number
  default     = 5
}

variable "eip_bandwidth_share_type" {
  description = "The share type of the EIP bandwidth"
  type        = string
  default     = "PER"

  validation {
    condition     = contains(["PER", "WHOLE"], var.eip_bandwidth_share_type)
    error_message = "The eip_bandwidth_share_type must be one of: PER, WHOLE."
  }
}

variable "eip_bandwidth_charge_mode" {
  description = "The charge mode of the EIP bandwidth"
  type        = string
  default     = "traffic"

  validation {
    condition     = contains(["traffic", "bandwidth"], var.eip_bandwidth_charge_mode)
    error_message = "The eip_bandwidth_charge_mode must be one of: traffic, bandwidth."
  }
}

# Variable definitions for SNAT Rule
variable "snat_source_type" {
  description = "The resource type of the SNAT rule"
  type        = number
  default     = 0
  nullable    = false

  validation {
    condition     = contains([0, 1], var.snat_source_type)
    error_message = "The snat_source_type must be one of: 0 (VPC scenario), 1 (Direct Connect scenario)."
  }
}

variable "snat_cidr" {
  description = "The CIDR block connected by SNAT rule (DC side, required when snat_source_type is 1)"
  type        = string
  default     = ""
  nullable    = false
}

variable "snat_description" {
  description = "The description of the SNAT rule"
  type        = string
  default     = ""
}

# Variable definitions for ECS Instance
variable "ecs_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "ecs_flavor_cpu_core_count" {
  description = "The CPU core count of the ECS instance flavor"
  type        = number
  default     = 2
}

variable "ecs_flavor_memory_size" {
  description = "The memory size of the ECS instance flavor"
  type        = number
  default     = 4
}

variable "ecs_flavor_id" {
  description = "The ID of the ECS instance flavor"
  type        = string
  default     = ""
  nullable    = false
}

variable "ecs_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "ecs_image_os" {
  description = "The OS of the ECS instance image"
  type        = string
  default     = "Ubuntu"
}

variable "ecs_security_group_name" {
  description = "The name of the security group for ECS instance"
  type        = string
  default     = "terraform-test-ecs-sg"
}

variable "ecs_instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "ecs_image_id" {
  description = "The ID of the ECS instance image"
  type        = string
  default     = ""
  nullable    = false
}

variable "ecs_system_disk_type" {
  description = "The type of the ECS instance system disk"
  type        = string
  default     = "SSD"
}

variable "ecs_system_disk_size" {
  description = "The size of the ECS instance system disk in GB"
  type        = number
  default     = 40
}

variable "ecs_admin_password" {
  description = "The password of the ECS instance administrator"
  type        = string
  sensitive   = true
}

variable "ecs_instance_tags" {
  description = "The tags of the ECS instance"
  type        = map(string)
  default     = {}
}
