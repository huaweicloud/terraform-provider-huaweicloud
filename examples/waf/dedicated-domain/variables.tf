# Variable definitions for authentication
variable "region_name" {
  description = "The region where the WAF dedicated domain is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for resources/data sources
variable "vpc_name" {
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "availability_zone" {
  description = "The availability zone to which the dedicated instance belongs"
  type        = string
  default     = ""
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "dedicated_instance_flavor_id" {
  description = "The flavor ID of the dedicated instance"
  type        = string
  default     = ""
}

variable "dedicated_instance_performance_type" {
  description = "The performance type of the dedicated instance"
  type        = string
  default     = "normal"
}

variable "dedicated_instance_cpu_core_count" {
  description = "The number of the dedicated instance CPU cores"
  type        = number
  default     = 4
}

variable "dedicated_instance_memory_size" {
  description = "The memory size of the dedicated instance"
  type        = number
  default     = 8
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "dedicated_instance_name" {
  description = "The WAF dedicated instance name"
  type        = string
}

variable "dedicated_instance_specification_code" {
  description = "The specification code of the dedicated instance"
  type        = string
  default     = "waf.instance.professional"
}

variable "policy_name" {
  description = "The WAF policy name"
  type        = string
}

variable "policy_level" {
  description = "The WAF policy level"
  type        = number
  default     = 1
}

variable "dedicated_mode_domain_name" {
  description = "The WAF dedicated mode domain name"
  type        = string
}

variable "dedicated_domain_client_protocol" {
  description = "The client protocol of the WAF dedicated domain"
  type        = string
  default     = "HTTP"
}

variable "dedicated_domain_server_protocol" {
  description = "The server protocol of the WAF dedicated domain"
  type        = string
  default     = "HTTP"
}

variable "dedicated_domain_address" {
  description = "The address of the WAF dedicated domain"
  type        = string
  default     = "192.168.0.14"
}

variable "dedicated_domain_port" {
  description = "The port of the WAF dedicated domain"
  type        = number
  default     = 8080
}

variable "dedicated_domain_type" {
  description = "The type of the WAF dedicated domain"
  type        = string
  default     = "ipv4"
}
