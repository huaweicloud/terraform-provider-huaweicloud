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

# Variable definitions for network resources
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

# Variable definitions for ECS instance
variable "ecs_instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "ecs_image_name" {
  description = "The name of the image used to create the ECS instance"
  type        = string
  default     = "Ubuntu 20.04 server 64bit"
}

variable "ecs_flavor_name" {
  description = "The flavor name of the ECS instance"
  type        = string
  default     = "s6.small.1"
}

variable "availability_zone" {
  description = "The availability zone where the ECS instance is located"
  type        = string
}

variable "ecs_tags" {
  description = "The tags of the ECS instance"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}

# Variable definitions for policy assignment
variable "policy_assignment_name" {
  description = "The name of the policy assignment"
  type        = string
}

variable "policy_assignment_description" {
  description = "The description of the policy assignment"
  type        = string
  default     = "Check if ECS instances have required tags"
}

variable "policy_definition_id" {
  description = "The ID of the policy definition"
  type        = string
}

variable "policy_assignment_policy_filter" {
  description = "The configuration used to filter resources"

  type = list(object({
    region            = string
    resource_provider = string
    resource_type     = string
    resource_id       = string
    tag_key           = string
    tag_value         = string
  }))

  default = []
}

variable "policy_assignment_parameters" {
  description = "The parameters of the policy assignment"
  type        = map(string)
  default     = {
    "tagKeys" = "[\"Owner\",\"Env\"]"
  }
}

variable "policy_assignment_tags" {
  description = "The tags of the policy assignment"
  type        = map(string)
  default     = {
    "Owner" = "terraform"
    "Env"   = "test"
  }
}
