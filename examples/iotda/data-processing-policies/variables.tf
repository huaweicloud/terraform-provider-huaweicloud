# Variable definitions for authentication
variable "region_name" {
  description = "The region where the IoTDA service is located"
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

# Variable definitions for the data flow control policy
variable "flow_control_policy_name" {
  description = "The name of the data flow control policy"
  type        = string
}

variable "flow_control_policy_description" {
  description = "The description of the data flow control policy"
  type        = string
  default     = ""
}

variable "flow_control_policy_limit" {
  description = "The flow control limit of the policy in tps"
  type        = number
}

# Variable definitions for the data backlog policy
variable "backlog_policy_name" {
  description = "The name of the data backlog policy"
  type        = string
}

variable "backlog_policy_description" {
  description = "The description of the data backlog policy"
  type        = string
  default     = ""
}

variable "backlog_policy_size" {
  description = "The size of data backlog in bytes"
  type        = string
}

variable "backlog_policy_time" {
  description = "The data backlog time in seconds"
  type        = string
}

variable "iotda_access_address" {
  description = "The HTTPS application access address of the IoTDA instance"
  type        = string
}
