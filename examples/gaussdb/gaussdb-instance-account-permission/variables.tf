# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GaussDB instance is located"
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
variable "vpc_name" {
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
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

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "security_group_name" {
  description = "The security group name"
  type        = string
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
  default     = "2379-2380,5000-5001,5432-5532,6000,6500,12016,20050"
}

variable "instance_password" {
  description = "The password for the GeminiDB Cassandra instance"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_name" {
  description = "The GeminiDB Cassandra instance name"
  type        = string
}

variable "instance_volume_type" {
  description = "The storage volume type"
  type        = string
  default     = "ULTRAHIGH"
}

variable "instance_volume_size" {
  description = "The storage volume size in GB"
  type        = number
  default     = 40
}

variable "db_name" {
  description = "The database name"
  type        = string
}

variable "db_owner" {
  description = "The database owner"
  type        = string
}

variable "database_account_name" {
  description = "The name of the database account"
  type        = string
}

variable "database_account_password" {
  description = "The password of the database account"
  type        = string
  default     = ""
  sensitive   = true
}

variable "is_login_only" {
  description = "Whether the database account supports login only. The valid values are true and false"
  type        = string
  default     = "false"
}

variable "schema_name" {
  description = "The schema name"
  type        = string
}

variable "permission_readonly" {
  description = "Whether the database account permission is read-only. The valid values are true and false"
  type        = string
  default     = "true"
}
