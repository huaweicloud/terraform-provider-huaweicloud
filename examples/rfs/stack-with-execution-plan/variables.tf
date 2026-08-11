# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RFS resources are located"
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
variable "stack_name" {
  description = "The name of the RFS resource stack"
  type        = string
}

variable "stack_description" {
  description = "The description of the RFS resource stack"
  type        = string
  default     = ""
}

variable "execution_plan_name" {
  description = "The name of the RFS execution plan"
  type        = string
}

variable "execution_plan_description" {
  description = "The description of the RFS execution plan"
  type        = string
  default     = ""
}

variable "execution_plan_template_body" {
  description = "The Terraform template body used by the RFS execution plan"
  type        = string
}
