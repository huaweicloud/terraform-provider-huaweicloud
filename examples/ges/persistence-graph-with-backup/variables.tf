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

variable "keep_backup" {
  description = "Whether to retain backups after the graph is deleted"
  type        = bool
  default     = true
}

variable "vertex_id_type" {
  description = "The vertex ID type for the persistence edition graph"
  type        = string
  default     = "hash"
}

variable "enable_encryption" {
  description = "Whether to enable data encryption for the graph instance"
  type        = bool
  default     = false
}

variable "kms_key_id" {
  description = "The KMS key ID for graph data encryption"
  type        = string
  default     = ""
  nullable    = false
}

variable "graph_tags" {
  description = "The tags of the GES graph instance"
  type        = map(string)
  default     = {}
}
