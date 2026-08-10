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
variable "function_name" {
  description = "The name of the FunctionGraph function"
  type        = string
}

variable "function_app" {
  description = "The group name of the FunctionGraph function"
  type        = string
  default     = "default"
}

variable "function_handler" {
  description = "The handler of the FunctionGraph function"
  type        = string
  default     = "index.handler"
}

variable "function_memory_size" {
  description = "The memory size of the FunctionGraph function in MB"
  type        = number
  default     = 128
}

variable "function_timeout" {
  description = "The timeout of the FunctionGraph function in seconds"
  type        = number
  default     = 3
}

variable "function_runtime" {
  description = "The runtime of the FunctionGraph function"
  type        = string
  default     = "Node.js12.13"
}

variable "function_code" {
  description = "The inline code content of the FunctionGraph function"
  type        = string
}

variable "private_provider_name" {
  description = "The name of the RFS private provider"
  type        = string
}

variable "private_provider_description" {
  description = "The description of the RFS private provider"
  type        = string
  default     = ""
}

variable "private_provider_version" {
  description = "The initial version number of the RFS private provider"
  type        = string
  default     = "1.0.0"
}

variable "private_provider_version_description" {
  description = "The initial version description of the RFS private provider"
  type        = string
  default     = ""
}

variable "provider_version_number" {
  description = "The version number of the RFS private provider version"
  type        = string
  default     = "2.0.0"
}

variable "provider_version_description" {
  description = "The description of the RFS private provider version"
  type        = string
  default     = ""
}
