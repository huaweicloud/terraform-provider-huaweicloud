# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RAM resource share is located"
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

# VPC variables
variable "vpc_name" {
  description = "The name of the VPC to be created and shared"
  type        = string
  default     = "ram-shared-vpc"
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

# Subnet variables
variable "subnet_name" {
  description = "The name of the subnet to be created and shared"
  type        = string
  default     = "ram-shared-subnet"
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet. If not specified, it will be automatically calculated from VPC CIDR"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet. If not specified, it will be automatically calculated"
  type        = string
  default     = ""
}

variable "availability_zone" {
  description = "The availability zone where the subnet will be created. If not specified, the first available zone will be used"
  type        = string
  default     = ""
}

# Security group variables
variable "secgroup_name" {
  description = "The name of the security group to be created and shared"
  type        = string
  default     = "ram-shared-secgroup"
}

variable "secgroup_description" {
  description = "The description of the security group"
  type        = string
  default     = "Security group for RAM resource sharing"
}

variable "secgroup_delete_default_rules" {
  description = "Whether to delete the default security group rules"
  type        = bool
  default     = false
}

# Resource share variables
variable "resource_urns" {
  description = "The list of URNs of one or more resources to be shared. If not specified, URNs will be automatically generated from created resources (VPC, subnet, security group)"
  type        = set(string)
  default     = []
}

variable "resource_share_name" {
  description = "The name of the resource share"
  type        = string
}

variable "description" {
  description = "The description of the resource share"
  type        = string
  default     = ""
}

variable "principals" {
  description = "The list of one or more principals (account IDs or organization IDs) to share resources with"
  type        = list(string)
}

variable "permission_ids" {
  description = "The list of RAM permissions associated with the resource share"
  type        = list(string)
  default     = []
}
