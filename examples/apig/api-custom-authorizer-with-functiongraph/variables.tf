# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
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

# Variable definitions for resources/data sources
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
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "function_name" {
  description = "The function name of the FunctionGraph"
  type        = string
}

variable "function_memory_size" {
  description = "The memory size (MB) for the function"
  type        = number
  default     = 128
}

variable "function_runtime" {
  description = "The runtime environment for the function"
  type        = string
  default     = "Python3.9"
}

variable "function_timeout" {
  description = "The timeout (seconds) for the function"
  type        = number
  default     = 3
}

variable "function_handler" {
  description = "The handler of the function"
  type        = string
  default     = "index.handler"
}

variable "function_code_type" {
  description = "The code type of the function"
  type        = string
  default     = "inline"
}

variable "function_app" {
  description = "The name of the application to which the function belongs"
  type        = string
  default     = "default"
}

variable "function_code" {
  description = "The code content of the FunctionGraph function"
  type        = string
}

variable "availability_zones" {
  description = "The availability zones to which the APIG instance belongs"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "instance_name" {
  description = "The instance name of the dedicated APIG"
  type        = string
}

variable "instance_edition" {
  description = "The edition of the APIG instance"
  type        = string
  default     = "BASIC"
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the APIG instance belongs"
  type        = string
  default     = null
}

variable "availability_zones_count" {
  description = "The number of availability zones to which the instance belongs"
  type        = number
  default     = 1
}

variable "custom_authorizer_name" {
  description = "The name of the custom authorizer"
  type        = string
}

variable "function_version" {
  description = "The version of the function"
  type        = string
  default     = "latest"
}

variable "custom_authorizer_type" {
  description = "The type of the custom authorizer"
  type        = string
  default     = "FRONTEND"
}

variable "custom_authorizer_network_type" {
  description = "The network type of the custom authorizer"
  type        = string
  default     = "V1"
}

variable "custom_authorizer_cache_age" {
  description = "The cache age of the custom authorizer"
  type        = number
  default     = 0
}

variable "custom_authorizer_is_body_send" {
  description = "Whether to send body in the custom authorizer"
  type        = bool
  default     = false
}

variable "custom_authorizer_use_data" {
  description = "The user data used to obtain backend access authorization"
  type        = string
  default     = null
}

variable "custom_authorizer_identity" {
  description = "The identity list of the custom authorizer"
  type        = list(object({
    name       = string
    location   = string
    validation = optional(string, null)
  }))
  default     = []
  nullable    = false
}

variable "response_name" {
  description = "The response name of the dedicated APIG"
  type        = string
}

variable "response_rules" {
  description = "The response rules of the dedicated APIG"
  type        = list(object({
    error_type  = string
    body        = string
    status_code = optional(number, null)
    headers     = optional(list(object({
      key   = string
      value = string
    })), [])
  }))
  default     = []
  nullable    = false
}

variable "group_name" {
  description = "The group name of the dedicated APIG"
  type        = string
}

variable "api_type" {
  description = "The type of the API"
  type        = string
  default     = "Public"
}

variable "api_name" {
  description = "The name of the API"
  type        = string
}

variable "api_request_protocol" {
  description = "The request protocol of the API"
  type        = string
  default     = "BOTH"
}

variable "api_request_method" {
  description = "The request method of the API"
  type        = string
  default     = "GET"
}

variable "api_request_path" {
  description = "The request path of the API"
  type        = string
}

variable "api_matching" {
  description = "The matching rule of the API"
  type        = string
  default     = "Exact"
}

variable "api_backend_params" {
  description = "The backend parameters of the API"
  type        = list(object({
    type              = string
    name              = string
    location          = string
    value             = string
    system_param_type = optional(string, null)
  }))
  nullable    = false
}

variable "api_func_graph_network_type" {
  description = "The network type of the FunctionGraph function"
  type        = string
  default     = "V1"
}

variable "api_func_graph_request_protocol" {
  description = "The request protocol of the FunctionGraph function"
  type        = string
  default     = "HTTPS"
}
