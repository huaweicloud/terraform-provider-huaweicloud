# Variable definitions for authentication
variable "region_name" {
  description = "The region where the NAT gateway and DNAT rule are located"
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

# Variable definitions for resources
variable "vpc_name" {
  description = "The VPC name"
  type        = string
  default     = "vpc-dnat-basic"
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "172.16.0.0/16"
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
  default     = "subnet-dnat-basic"
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = "172.16.10.0/24"
}

variable "subnet_gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = true
}

variable "nat_gateway_name" {
  description = "The NAT gateway name"
  type        = string
  default     = "nat-gateway-dnat-basic"
}

variable "nat_gateway_description" {
  description = "The description of the NAT gateway"
  type        = string
  default     = ""
}

variable "nat_gateway_spec" {
  description = "The specification of the NAT gateway"
  type        = string
  default     = "1"

  validation {
    condition     = contains(["1", "2", "3", "4"], var.nat_gateway_spec)
    error_message = "The nat_gateway_spec must be one of: 1, 2, 3, 4."
  }
}

variable "eip_bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
  default     = "eip-dnat-basic"
}

variable "eip_bandwidth_size" {
  description = "The size of the EIP bandwidth (Mbit/s)"
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

variable "frontend_protocol" {
  description = "The protocol used between the client and NAT gateway"
  type        = string
  default     = "tcp"

  validation {
    condition     = contains(["tcp", "udp", "any"], var.frontend_protocol)
    error_message = "The frontend_protocol must be one of: tcp, udp, any."
  }
}

variable "backend_port" {
  description = "The port on the backend ECS instance that receives DNAT traffic"
  type        = number
  default     = 22
}

variable "frontend_port" {
  description = "The port on the public EIP that clients use to access the DNAT service"
  type        = number
  default     = 22
}

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
  description = "The flavor ID of the backend ECS instance"
  type        = string
  default     = ""
  nullable    = true
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

variable "security_group_name" {
  description = "The security group name of the backend instance"
  type        = string
  default     = "sg-dnat-backend"
}

variable "backend_protocol" {
  description = "The protocol used between the NAT gateway and backend ECS instance"
  type        = string
  default     = "tcp"
}

variable "ingress_cidr" {
  description = "The CIDR block that is allowed to access the DNAT service from the Internet"
  type        = string
  default     = "0.0.0.0/0"
}

variable "instance_name" {
  description = "The name of the backend ECS instance"
  type        = string
  default     = "ecs-dnat-backend"
}

variable "ecs_image_id" {
  description = "The image ID of the backend ECS instance"
  type        = string
  default     = ""
  nullable    = true
}

variable "ecs_system_disk_type" {
  description = "The system disk type of the ECS instance"
  type        = string
  default     = "SSD"
}

variable "ecs_system_disk_size" {
  description = "The system disk size of the ECS instance (GB)"
  type        = number
  default     = 40
}

variable "ecs_admin_password" {
  description = "The administrator password of the ECS instance"
  type        = string
  sensitive   = true
  default     = ""
}

variable "ecs_instance_tags" {
  description = "The tags of the ECS instance"
  type        = map(string)
  default     = {}
}
