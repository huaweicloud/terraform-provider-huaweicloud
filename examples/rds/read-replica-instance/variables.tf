# Variable definitions for authentication
variable "region_name" {
  description = "The region where the PostgreSQL RDS instance is located"
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
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "availability_zones" {
  description = "The list of availability zones to which the RDS instance belong"
  type        = list(string)
  default     = []
  nullable    = false

  validation {
    condition     = length(var.availability_zones) == 0 || length(var.availability_zones) >= 2
    error_message = "The availability zones must be a list of strings"
  }
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

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "instance_flavors_filter" {
  description = "The filter configuration of the RDS MySQL flavor instance mode for the RDS instance"
  type        = list(object({
    db_type       = optional(string, "MySQL")
    db_version    = optional(string, "8.0")
    instance_mode = optional(string, "ha")
    group_type    = optional(string, "general")
    vcpus         = optional(number, 4)
    memory        = optional(number, 8)
  }))

  default = [
    {
      instance_mode = "ha"
    },
    {
      instance_mode = "replica"
    }
  ]

  validation {
    condition = length(var.instance_flavors_filter) == 2 && length(setintersection([for o in var.instance_flavors_filter : lookup(o, "instance_mode", "")], ["ha", "replica"])) == 2
    error_message = "The instance_flavors_filter must contain at least two elements and must have both 'ha' and 'replica'."
  }
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
}

variable "instance_db_port" {
  description = "The database port"
  type        = number
  default     = 5432
}

variable "instance_password" {
  description = "The password for the RDS instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The name of the RDS instance"
  type        = string
}

variable "instance_flavor_id" {
  description = "The flavor ID of the RDS instance"
  type        = string
  default     = ""
}

variable "ha_replication_mode" {
  description = "The HA replication mode of the RDS instance"
  type        = string
  default     = "async"
}

variable "instance_volume_type" {
  description = "The storage volume type"
  type        = string
  default     = "CLOUDSSD"
}

variable "instance_volume_size" {
  description = "The storage volume size in GB"
  type        = number
  default     = 40
}

variable "instance_backup_time_window" {
  description = "The backup time window in HH:MM-HH:MM format"
  type        = string
}

variable "instance_backup_keep_days" {
  description = "The number of days to retain backups"
  type        = number
}

variable "replica_instance_name" {
  description = "The name of the read replica instance"
  type        = string
}

variable "replica_instance_flavor_id" {
  description = "The flavor ID of the read replica instance"
  type        = string
  default     = ""
}

variable "replica_instance_volume_type" {
  description = "The storage volume type"
  type        = string
  default     = "CLOUDSSD"
}

variable "replica_instance_volume_size" {
  description = "The storage volume size in GB"
  type        = number
  default     = 40
}
