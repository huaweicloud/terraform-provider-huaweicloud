# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GES backup is located"
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

# Variable definitions for huaweicloud_vpc
variable "vpc_name" {
  description = "The VPC name for the GES graph"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
}

# Variable definitions for huaweicloud_vpc_subnet
variable "subnet_name" {
  description = "The subnet name for the GES graph"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
}

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
}

# Variable definitions for huaweicloud_networking_secgroup
variable "security_group_name" {
  description = "The security group name for the GES graph"
  type        = string
}

# Variable definitions for huaweicloud_ges_graph
variable "graph_name" {
  description = "The GES graph name"
  type        = string
}

variable "graph_size_type_index" {
  description = "The graph size type index"
  type        = string
  default     = "1"
}

variable "graph_cpu_arch" {
  description = "The CPU architecture type of the GES graph"
  type        = string
  default     = "x86_64"
}

variable "graph_crypt_algorithm" {
  description = "The cryptography algorithm of the GES graph"
  type        = string
}

variable "graph_enable_https" {
  description = "Whether to enable HTTPS for the GES graph"
  type        = bool
  default     = false
}

variable "graph_tags" {
  description = "The key/value pairs to associate with the GES graph"
  type        = map(string)
  default     = {
    key = "val"
    foo = "bar"
  }
  nullable    = false
}
