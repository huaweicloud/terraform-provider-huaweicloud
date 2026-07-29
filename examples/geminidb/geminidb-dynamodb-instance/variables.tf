# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GeminiDB DynamoDB instance is located"
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

# Variable definitions for VPC and subnet
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
}

variable "subnet_name" {
  description = "The name of the VPC subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the VPC subnet, defaults to a /24 derived from the VPC CIDR"
  type        = string
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the VPC subnet, defaults to the first IP of the derived subnet CIDR"
  type        = string
  default     = ""
}

variable "vcpus" {
  description = "The number of vCPUs of the flavor to query"
  type        = number
  default     = 2
}

variable "availability_zone" {
  description = "The availability zone of the GeminiDB DynamoDB instance, uses the first AZ from data source if empty"
  type        = string
  default     = ""
}

# Variable definitions for security group
variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_password" {
  description = "The password of the GeminiDB DynamoDB instance, generated randomly if left empty"
  type        = string
  sensitive   = true
  default     = ""
}

# Variable definitions for GeminiDB DynamoDB instance
variable "instance_name" {
  description = "The name of the GeminiDB DynamoDB instance"
  type        = string
}

variable "instance_mode" {
  description = "The instance type of the GeminiDB DynamoDB instance"
  type        = string
  default     = "Cluster"
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the GeminiDB DynamoDB instance"
  type        = string
  default     = "0"
}

variable "instance_ssl_option" {
  description = "Whether SSL is enabled for the GeminiDB DynamoDB instance"
  type        = string
  default     = "on"
}

variable "instance_flavor_num" {
  description = "The node quantity of the GeminiDB DynamoDB instance"
  type        = number
  default     = 3
}

variable "instance_flavor_size" {
  description = "The disk size (GB) of the GeminiDB DynamoDB instance"
  type        = number
  default     = 200
}

variable "instance_flavor_storage" {
  description = "The disk type of the GeminiDB DynamoDB instance"
  type        = string
  default     = "ULTRAHIGH"
}

variable "instance_flavor_spec_code" {
  description = "The resource specification code, uses the first flavor from data source if empty"
  type        = string
  default     = ""
}

variable "instance_backup_time_window" {
  description = "The backup time window of the GeminiDB DynamoDB instance"
  type        = string
  default     = "03:00-04:00"
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backup files of the GeminiDB DynamoDB instance"
  type        = number
  default     = 14
}

variable "charging_mode" {
  description = "The charging mode of the GeminiDB DynamoDB instance, postPaid or prePaid"
  type        = string
  default     = "prePaid"
}

variable "period_unit" {
  description = "The charging period unit of the GeminiDB DynamoDB instance, month or year"
  type        = string
  default     = "month"
}

variable "auto_renew" {
  description = "Whether auto-renew is enabled for the GeminiDB DynamoDB instance"
  type        = string
  default     = "true"
}

variable "period" {
  description = "The charging period of the GeminiDB DynamoDB instance"
  type        = number
  default     = 1
}

variable "tags" {
  description = "The key/value pairs to associate with the GeminiDB DynamoDB instance"
  type        = map(string)
  default     = {
    foo = "bar"
    key = "value"
  }
}

# Variable definitions for GeminiDB backup
variable "backup_name" {
  description = "The name of the GeminiDB backup"
  type        = string
  default     = "tf_test_backup"
}

variable "backup_description" {
  description = "The description of the GeminiDB backup"
  type        = string
  default     = "Created by Terraform"
}
