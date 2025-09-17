# Variables definitions for authorization
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

# Variables definitions for resource/data source
variable "instance_flavor_id" {
  description = "The flavor ID of the AS instance"
  type        = string
  default     = ""
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the AS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the AS instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the AS instance flavor"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the AS instance"
  type        = string
  default     = ""
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
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""

  validation {
    condition     = (var.subnet_cidr != "" && var.subnet_gateway_ip != "") || (var.subnet_cidr == "" && var.subnet_gateway_ip == "")
    error_message = "The 'subnet_cidr' and 'subnet_gateway_ip' is not allowed for only one of them to be empty"
  }
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "keypair_name" {
  description = "The name of the key pair that is used to access the AS instance"
  type        = string
}

variable "keypair_public_key" {
  description = "The public key of the key pair that is used to access the AS instance"
  type        = string
  default     = ""
}

variable "configuration_name" {
  description = "The name of the AS configuration"
  type        = string
}

variable "disk_configurations" {
  description = "The disk configurations for the AS instance"
  type        = list(object({
    disk_type   = string
    volume_type = string
    volume_size = number
  }))

  nullable = false

  validation {
    condition     = length(var.disk_configurations) > 0 && length([for v in var.disk_configurations : v if v.disk_type == "SYS"]) == 1
    error_message = "The 'disk_configurations' is not allowed to be empty and only one system disk is allowed"
  }
}

variable "group_name" {
  description = "The name of the AS group"
  type        = string
}

variable "desire_instance_number" {
  description = "The desired number of scaling instances in the AS group"
  type        = number
  default     = 2
}

variable "min_instance_number" {
  description = "The minimum number of scaling instances in the AS group"
  type        = number
  default     = 0
}

variable "max_instance_number" {
  description = "The maximum number of scaling instances in the AS group"
  type        = number
  default     = 10
}

variable "is_delete_publicip" {
  description = "Whether to delete the public IP address of the scaling instances when the AS group is deleted"
  type        = bool
  default     = true
}

variable "is_delete_instances" {
  description = "Whether to delete the scaling instances when the AS group is deleted"
  type        = bool
  default     = true
}

variable "topic_name" {
  description = "The name of the SMN topic"
  type        = string
}

variable "alarm_rule_name" {
  description = "The name of the CES alarm rule"
  type        = string
}

variable "rule_conditions" {
  description = "The conditions of the alarm rule"
  type        = list(object({
    alarm_level         = optional(number, 2)
    metric_name         = string
    period              = number
    filter              = string
    comparison_operator = string
    suppress_duration   = optional(number, 0)
    value               = number
    count               = number
  }))

  nullable = false

  validation {
    condition     = length(var.rule_conditions) > 0
    error_message = "The 'rule_conditions' is not allowed to be empty"
  }
}

variable "scaling_up_policy_name" {
  description = "The name of the scaling up policy"
  type        = string
}

variable "scaling_up_cool_down_time" {
  description = "The cool down time of the scaling up policy"
  type        = number
  default     = 300
}

variable "scaling_up_instance_number" {
  description = "The number of instances to add when the scaling up policy is triggered"
  type        = number
  default     = 1
}

variable "scaling_down_policy_name" {
  description = "The name of the scaling down policy"
  type        = string
}

variable "scaling_down_cool_down_time" {
  description = "The cool down time of the scaling down policy"
  type        = number
  default     = 300
}

variable "scaling_down_instance_number" {
  description = "The number of instances to remove when the scaling down policy is triggered"
  type        = number
  default     = 1
}
