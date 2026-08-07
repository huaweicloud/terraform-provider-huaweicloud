# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GES resources are located"
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
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "graph_name" {
  description = "The name of the GES graph instance"
  type        = string
}

variable "graph_size_type_index" {
  description = "The graph size type index"
  type        = string
  default     = "3"
}

variable "enable_https" {
  description = "Whether to enable the security mode (HTTPS). This may impact performance"
  type        = bool
  default     = false
}

variable "enable_multi_az" {
  description = "Whether to enable cross-AZ deployment for high availability"
  type        = bool
  default     = true
}

variable "enable_audit" {
  description = "Whether to enable graph audit logs via LTS"
  type        = bool
  default     = false
}

variable "audit_log_group_name" {
  description = "The LTS log group name for graph audit logs"
  type        = string
  default     = ""
  nullable    = false
}

variable "graph_tags" {
  description = "The tags of the GES graph instance"
  type        = map(string)
  default     = {}
}
