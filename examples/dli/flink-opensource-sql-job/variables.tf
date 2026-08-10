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
  description = "The CIDR block of the elastic resource pool."
  type        = string
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
  nullable    = false
}

variable "elastic_resource_pool_label" {
  description = "The label of the elastic resource pool"
  type        = map(string)
  default     = {}
}

variable "queue_name" {
  description = "The name of the DLI general queue used to run the Flink OpenSource SQL job"
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

variable "job_name" {
  description = "The name of the Flink OpenSource SQL job"
  type        = string
}

variable "job_sql" {
  description = "The Flink SQL statement that includes source, query and sink"
  type        = string
}

variable "job_flink_version" {
  description = "The Flink version"
  type        = string
}

variable "job_execution_agency_urn" {
  description = "The agency URN authorized to DLI"
  type        = string
  default     = ""
}

variable "job_cu_number" {
  description = "The number of CUs selected for the Flink OpenSource SQL job"
  type        = number
  default     = 2
}

variable "job_parallel_number" {
  description = "The parallelism of the Flink OpenSource SQL job"
  type        = number
  default     = 1
}

variable "job_manager_cu_number" {
  description = "The number of CUs in the JobManager"
  type        = number
  default     = 1
}

variable "job_tm_cus" {
  description = "The number of CUs for each Task Manager"
  type        = number
  default     = 1
}

variable "job_tm_slot_num" {
  description = "The number of slots in each Task Manager"
  type        = number
  default     = null
}

variable "job_checkpoint_enabled" {
  description = "Whether to enable the automatic job snapshot function"
  type        = bool
  default     = false
}

variable "job_checkpoint_mode" {
  description = "The snapshot mode."
  type        = string
  default     = "exactly_once"
}

variable "job_checkpoint_interval" {
  description = "The snapshot interval in seconds"
  type        = number
  default     = 10
}

variable "job_obs_bucket" {
  description = "The OBS bucket used to save snapshots or logs"
  type        = string
  default     = ""
}

variable "job_log_enabled" {
  description = "Whether to enable uploading job logs to the OBS bucket"
  type        = bool
  default     = false
}

variable "job_smn_topic" {
  description = "The SMN topic used to receive job failure notifications"
  type        = string
  default     = ""
}

variable "job_restart_when_exception" {
  description = "Whether to enable the function of automatically restarting a job upon job exceptions"
  type        = bool
  default     = false
}

variable "job_resume_max_num" {
  description = "The maximum number of retry times upon exceptions"
  type        = number
  default     = -1
}

variable "job_idle_state_retention" {
  description = "The idle state retention in seconds"
  type        = number
  default     = 1
}

variable "job_runtime_config" {
  description = "The custom optimization parameters when the Flink job is running"
  type        = map(string)
  default     = {}
}

variable "job_tags" {
  description = "The key/value pairs to associate with the Flink OpenSource SQL job"
  type        = map(string)
  default     = {}
}

variable "job_description" {
  description = "The description of the Flink OpenSource SQL job"
  type        = string
  default     = ""
}
