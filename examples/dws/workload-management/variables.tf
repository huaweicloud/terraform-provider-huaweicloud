# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DWS resources are located"
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
  description = "The availability zone of the DWS cluster"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
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

variable "security_group_delete_default_rules" {
  description = "Whether to delete the default rules of the security group"
  type        = bool
  default     = true
}

variable "security_group_rule_ports" {
  description = "The security group ingress rule ports for DWS"
  type        = string
  default     = "8000-10000"
}

variable "cluster_node_type" {
  description = "The flavor of the DWS cluster node"
  type        = string
  default     = ""
  nullable    = false
}

variable "cluster_version" {
  description = "The version of the DWS cluster"
  type        = string
  default     = ""
  nullable    = false
}

variable "cluster_vcpus" {
  description = "The vcpus of the DWS cluster"
  type        = number
  default     = 4
}

variable "cluster_memory" {
  description = "The memory of the DWS cluster"
  type        = number
  default     = 32
}

variable "cluster_datastore_type" {
  description = "The datastore type of the DWS cluster"
  type        = string
  default     = "dws"
}

variable "cluster_name" {
  description = "The name of the DWS cluster"
  type        = string
}

variable "cluster_number_of_node" {
  description = "The number of nodes in the DWS cluster"
  type        = number
  default     = 3
}

variable "cluster_number_of_cn" {
  description = "The number of CN nodes in the DWS cluster"
  type        = number
  default     = 3
}

variable "cluster_admin_user_name" {
  description = "The administrator username of the DWS cluster"
  type        = string
}

variable "cluster_admin_user_pwd" {
  description = "The administrator password of the DWS cluster"
  type        = string
  sensitive   = true
}

variable "cluster_volume_type" {
  description = "The volume type of the DWS cluster"
  type        = string
  default     = "SSD"
}

variable "cluster_volume_capacity" {
  description = "The volume capacity of the DWS cluster in GB"
  type        = string
  default     = "100"
}

variable "workload_queue_name" {
  description = "The name of the workload queue"
  type        = string
}

variable "workload_queue_configurations" {
  description = "The configurations of the workload queue"

  type = list(object({
    resource_name  = string
    resource_value = string
  }))
}

variable "user_name" {
  description = "The name of the cluster user"
  type        = string
}

variable "user_password" {
  description = "The password of the cluster user"
  type        = string
  sensitive   = true
}

variable "user_description" {
  description = "The description of the cluster user"
  type        = string
  default     = ""
  nullable    = false
}

variable "user_cascade" {
  description = "Whether to cascade delete dependencies when deleting the user or role"
  type        = bool
  default     = true
}

variable "user_login" {
  description = "Whether to allow the user to log in"
  type        = bool
  default     = true
}

variable "user_create_role" {
  description = "Whether to grant the permission to create roles"
  type        = bool
  default     = true
}

variable "user_create_db" {
  description = "Whether to grant the permission to create databases"
  type        = bool
  default     = true
}

variable "user_system_admin" {
  description = "Whether to grant the system administrator permission"
  type        = bool
  default     = null
}

variable "user_audit_admin" {
  description = "Whether to grant the audit administrator permission"
  type        = bool
  default     = null
}

variable "user_inherit" {
  description = "Whether to inherit permissions from roles"
  type        = bool
  default     = true
}

variable "user_use_ft" {
  description = "Whether to grant the external table permission"
  type        = bool
  default     = null
}

variable "user_conn_limit" {
  description = "The maximum number of concurrent connections. -1 means unlimited"
  type        = number
  default     = -1
}

variable "user_replication" {
  description = "Whether to grant the replication permission"
  type        = bool
  default     = null
}

variable "user_valid_begin" {
  description = "The valid begin time of the cluster user"
  type        = string
  default     = ""
  nullable    = false
}

variable "user_valid_until" {
  description = "The valid until time of the cluster user"
  type        = string
  default     = ""
  nullable    = false
}

variable "user_grant_list" {
  description = "The set of grants for the cluster user"

  type = list(object({
    type                 = string
    database             = optional(string)
    schema_name          = optional(string)
    object_name          = optional(string)
    all_object           = optional(bool)
    future               = optional(bool)
    future_object_owners = optional(string)
    column_names         = optional(list(string))

    privileges = list(object({
      permission = string
      grant_with = bool
    }))
  }))

  default = []
}

variable "workload_plan_name" {
  description = "The name of the workload plan"
  type        = string
}

variable "workload_plan_stage_name" {
  description = "The name of the workload plan stage"
  type        = string
}

variable "workload_plan_stage_month" {
  description = "The month of the workload plan stage"
  type        = string
  default     = null
}

variable "workload_plan_stage_day" {
  description = "The day of the workload plan stage"
  type        = string
  default     = null
}

variable "workload_plan_stage_start_time" {
  description = "The start time of the workload plan stage. The format is hh:mm:ss"
  type        = string
  default     = "00:00:00"
}

variable "workload_plan_stage_end_time" {
  description = "The end time of the workload plan stage. The format is hh:mm:ss"
  type        = string
  default     = "23:59:59"
}

variable "workload_plan_stage_configurations" {
  description = "The configurations of the workload plan stage"

  type = list(object({
    resource_name        = string
    resource_value       = string
    value_unit           = optional(string)
    resource_description = optional(string)
  }))
}

variable "exception_rule_name" {
  description = "The name of the cluster exception rule"
  type        = string
}

variable "exception_rule_configurations" {
  description = "The configurations of the exception rule"

  type = list(object({
    key   = string
    value = string
  }))
}
