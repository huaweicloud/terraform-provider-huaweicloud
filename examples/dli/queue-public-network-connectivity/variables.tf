# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DLI resources are located"
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
variable "elastic_resource_pool_name" {
  description = "The name of the DLI elastic resource pool"
  type        = string
}

variable "elastic_resource_pool_description" {
  description = "The description of the elastic resource pool"
  type        = string
  default     = ""
}

variable "elastic_resource_pool_min_cu" {
  description = "The minimum number of CUs for the elastic resource pool"
  type        = number
  default     = 16
}

variable "elastic_resource_pool_max_cu" {
  description = "The maximum number of CUs for the elastic resource pool"
  type        = number
  default     = 64
}

variable "elastic_resource_pool_cidr" {
  description = "The CIDR block of the elastic resource pool. This CIDR must not overlap with the VPC CIDR"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
  nullable    = false
}

variable "queue_name" {
  description = "The name of the DLI exclusive queue"
  type        = string
}

variable "queue_type" {
  description = "The type of the DLI queue. The valid values are sql and general"
  type        = string
  default     = "sql"

  validation {
    condition     = contains(["sql", "general"], var.queue_type)
    error_message = "The queue_type valid value must be `sql` or `general`."
  }
}

variable "queue_cu_count" {
  description = "The CU count of the DLI queue"
  type        = number
  default     = 16
}

variable "queue_description" {
  description = "The description of the DLI queue"
  type        = string
  default     = ""
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet. If empty, it is calculated from the VPC CIDR"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet. If empty, it is calculated from the subnet CIDR"
  type        = string
  default     = ""
  nullable    = false
}

variable "datasource_connection_name" {
  description = "The name of the DLI enhanced datasource connection"
  type        = string
}

variable "datasource_connection_hosts" {
  description = "The list of custom hosts for the enhanced datasource connection"

  type = list(object({
    name = string
    ip   = string
  }))

  default  = []
  nullable = false
}

variable "datasource_connection_routes" {
  description = "The list of custom routes for the enhanced datasource connection. Each cidr should be the public destination network to access"

  type = list(object({
    name = string
    cidr = string
  }))
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "eip_bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
}

variable "eip_bandwidth_size" {
  description = "The size of the EIP bandwidth in Mbps"
  type        = number
  default     = 5
}

variable "eip_bandwidth_share_type" {
  description = "The share type of the EIP bandwidth"
  type        = string
  default     = "PER"
}

variable "eip_bandwidth_charge_mode" {
  description = "The charge mode of the EIP bandwidth"
  type        = string
  default     = "traffic"
}

variable "nat_gateway_name" {
  description = "The name of the NAT gateway"
  type        = string
}

variable "nat_gateway_spec" {
  description = "The specification of the NAT gateway."
  type        = string
  default     = "1"

  validation {
    condition     = contains(["1", "2", "3", "4"], var.nat_gateway_spec)
    error_message = "The nat_gateway_spec valid value must be `1`, `2`, `3` or `4`."
  }
}

variable "nat_gateway_description" {
  description = "The description of the NAT gateway"
  type        = string
  default     = ""
}

variable "snat_description" {
  description = "The description of the SNAT rule"
  type        = string
  default     = ""
}
