# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GaussDB instance is located"
  type        = string
  nullable    = false
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

# Variable definitions for network resources
variable "vpc_name" {
  description = "The VPC name"
  type        = string
  nullable    = false
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  nullable    = false
  default     = "172.16.0.0/16"
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
  nullable    = false
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  nullable    = false
  default     = ""
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  nullable    = false
  default     = ""
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
  nullable    = false
}

# Ports for GaussDB instances (all required, must be opened in security group):
# Required (system reserved ports, must be opened for instance to run):
# 2379-2380: ETCD service ports, used for distributed configuration sharing and service discovery
# 5000-5001: CMS (Cluster Management Server) ports, used for managing instance status
# 6000:      GTM (Global Transaction Manager) port, used for managing transaction status
# 6500:      Internal service port
# 12016:     Internal communication port
# 20050:     Internal service port
# Required (must be adjusted based on instance_db_port):
# 5432-5532: Database port range，database port 5432 and its subsequent 100 ports
variable "security_group_rule_ports" {
  description = "The security group ingress rule ports"
  type        = string
  nullable    = false
  default     = "2379-2380,5000-5001,5432-5532,6000,6500,12016,20050"
}

# Variable definitions for GaussDB instance
# If not specified, the first 3 available zones will be used automatically.
variable "instance_availability_zones" {
  description = "The availability zones for the GaussDB instance, separated by commas"
  type        = string
  nullable    = false
  default     = ""
}

variable "instance_name" {
  description = "The name of the GaussDB instance"
  type        = string
  nullable    = false
}

variable "instance_flavor" {
  description = "The flavor of the GaussDB instance"
  type        = string
  nullable    = false
  default     = "gaussdb.opengauss.ee.c3.xlarge.x864.ha"
}

variable "instance_password" {
  description = "The password for the GaussDB instance"
  type        = string
  sensitive   = true
  nullable    = false
  default     = ""
}

variable "instance_db_port" {
  description = "The database port of the GaussDB instance"
  type        = number
  nullable    = false
  default     = 5432
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the GaussDB instance"
  type        = string
  default     = null
}

variable "instance_ha_mode" {
  description = "The HA mode of the GaussDB instance"
  type        = string
  nullable    = false
  default     = "centralization_standard"
}

variable "instance_ha_replication_mode" {
  description = "The HA replication mode of the GaussDB instance"
  type        = string
  nullable    = false
  default     = "sync"
}

variable "instance_ha_consistency" {
  description = "The HA consistency of the GaussDB instance"
  type        = string
  nullable    = false
  default     = "strong"
}

variable "instance_volume_type" {
  description = "The storage volume type of the GaussDB instance"
  type        = string
  nullable    = false
  default     = "ULTRAHIGH"
}

variable "instance_volume_size" {
  description = "The storage volume size (GB) of the GaussDB instance"
  type        = number
  nullable    = false
  default     = 40
}

# Variable definitions for manual backup
variable "backup_name" {
  description = "The name for the manual backup"
  type        = string
  nullable    = false
}

variable "backup_description" {
  description = "The description for the manual backup"
  type        = string
  nullable    = false
  default     = ""
}
