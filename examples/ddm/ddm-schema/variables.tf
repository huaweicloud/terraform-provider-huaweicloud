# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DDM instance is located"
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
variable "availability_zones" {
  description = "The availability zones to which the DDM instance belongs"
  type        = list(string)
  default     = []
  nullable    = false
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

variable "rds_instance_password" {
  description = "The password of the RDS instance"
  type        = string
  sensitive   = true
  default     = ""
  nullable    = false
}

variable "instance_flavor" {
  description = "The flavor of the RDS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "database_type" {
  description = "The database type of the RDS instance"
  type        = string
  default     = "MySQL"
}

variable "database_version" {
  description = "The database version of the RDS instance"
  type        = string
  default     = "5.7"
}

variable "instance_mode" {
  description = "The mode of the RDS instance"
  type        = string
  default     = "single"
}

variable "instance_group_type" {
  description = "The performance specification"
  type        = string
  default     = "dedicated"
}

variable "instance_flavor_vcpus" {
  description = "The number of vCPUs for the RDS instance flavor"
  type        = number
  default     = 2
}

variable "rds_instance_name" {
  description = "The name of the RDS instance"
  type        = string
}

variable "database_port" {
  description = "The port of the RDS instance"
  type        = number
  default     = 3306
}

variable "volume_type" {
  description = "The volume type of the RDS instance"
  type        = string
  default     = "CLOUDSSD"
}

variable "volume_size" {
  description = "The volume size of the RDS instance"
  type        = number
  default     = 40
}

variable "instance_engine_id" {
  description = "The engine ID of the DDM instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the DDM instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "ddm_instance_name" {
  description = "The name of the DDM instance"
  type        = string
}

variable "instance_node_num" {
  description = "The number of nodes in the DDM instance"
  type        = number
  default     = 2
}

variable "instance_parameters" {
  description = "The parameters of the DDM instance"

  type = list(object({
    name  = string
    value = string
  }))

  default = []
}

variable "schema_name" {
  description = "The name of the DDM schema"
  type        = string
}

variable "schema_shard_mode" {
  description = "The shard mode of the DDM schema"
  type        = string
  default     = "single"
}

variable "schema_shard_number" {
  description = "The number of shards in the same working mode"
  type        = number
  default     = 1
}
