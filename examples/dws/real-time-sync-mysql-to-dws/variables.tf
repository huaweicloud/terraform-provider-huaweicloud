# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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
variable "availability_zone" {
  description = "The availability zone. If empty, the first available zone is used"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_name" {
  description = "The name of the VPC shared by RDS and DWS"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC. Must differ from the DLI elastic resource pool CIDR"
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
  nullable    = false
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

variable "security_group_name" {
  description = "The name of the security group shared by the RDS instance and DWS cluster"
  type        = string
}

variable "security_group_delete_default_rules" {
  description = "Whether to delete the default rules of the security group"
  type        = bool
  default     = true
}

variable "dws_port" {
  description = "The service port of the DWS cluster"
  type        = number
  default     = 8000
}

variable "rds_db_port" {
  description = "The database port of the RDS MySQL instance"
  type        = number
  default     = 3306
}

variable "elastic_resource_pool_cidr" {
  description = "The CIDR block of the DLI elastic resource pool. Must differ from the VPC CIDR"
  type        = string
}

variable "rds_flavor_id" {
  description = "The flavor ID of the RDS instance. If empty, it is queried from huaweicloud_rds_flavors"
  type        = string
  default     = ""
  nullable    = false
}

variable "rds_db_version" {
  description = "The MySQL version of the RDS instance"
  type        = string
  default     = "5.7"
}

variable "rds_instance_mode" {
  description = "The instance mode used to query RDS flavors"
  type        = string
  default     = "single"
}

variable "rds_flavor_vcpus" {
  description = "The vCPUs used to query RDS flavors"
  type        = number
  default     = 2
}

variable "rds_instance_name" {
  description = "The name of the RDS MySQL instance"
  type        = string
}

variable "rds_db_password" {
  description = "The root password of the RDS MySQL instance"
  type        = string
  sensitive   = true
}

variable "rds_volume_type" {
  description = "The volume type of the RDS instance"
  type        = string
  default     = "CLOUDSSD"
}

variable "rds_volume_size" {
  description = "The volume size of the RDS instance in GB"
  type        = number
  default     = 40
}

variable "dws_node_type" {
  description = "The flavor of the DWS cluster node. If empty, it is queried from huaweicloud_dws_flavors"
  type        = string
  default     = ""
  nullable    = false
}

variable "dws_version" {
  description = "The version of the DWS cluster. If empty, it is queried from huaweicloud_dws_flavors"
  type        = string
  default     = ""
  nullable    = false
}

variable "dws_flavor_vcpus" {
  description = "The vCPUs used to query DWS flavors"
  type        = number
  default     = 4
}

variable "dws_flavor_memory" {
  description = "The memory used to query DWS flavors"
  type        = number
  default     = 32
}

variable "dws_datastore_type" {
  description = "The datastore type of the DWS cluster"
  type        = string
  default     = "dws"
}

variable "dws_cluster_name" {
  description = "The name of the DWS cluster"
  type        = string
}

variable "dws_number_of_node" {
  description = "The number of nodes in the DWS cluster"
  type        = number
  default     = 3
}

variable "dws_number_of_cn" {
  description = "The number of CN nodes in the DWS cluster"
  type        = number
  default     = 3
}

variable "dws_admin_user_name" {
  description = "The administrator username of the DWS cluster"
  type        = string
  default     = "dbadmin"
}

variable "dws_admin_user_pwd" {
  description = "The administrator password of the DWS cluster"
  type        = string
  sensitive   = true
}

variable "dws_volume_type" {
  description = "The volume type of the DWS cluster"
  type        = string
  default     = "SSD"
}

variable "dws_volume_capacity" {
  description = "The volume capacity of the DWS cluster in GB"
  type        = string
  default     = "100"
}

variable "elastic_resource_pool_name" {
  description = "The name of the DLI elastic resource pool"
  type        = string
}

variable "elastic_resource_pool_description" {
  description = "The description of the DLI elastic resource pool"
  type        = string
  default     = ""
}

variable "elastic_resource_pool_min_cu" {
  description = "The minimum number of CUs for the DLI elastic resource pool"
  type        = number
  default     = 16
}

variable "elastic_resource_pool_max_cu" {
  description = "The maximum number of CUs for the DLI elastic resource pool"
  type        = number
  default     = 64
}

variable "elastic_resource_pool_label" {
  description = "The label of the DLI elastic resource pool"
  type        = map(string)

  default = {
    spec = "basic"
  }
}

variable "queue_name" {
  description = "The name of the DLI general queue"
  type        = string
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

variable "datasource_connection_name" {
  description = "The name of the DLI enhanced datasource connection"
  type        = string
}
