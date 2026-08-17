# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DRS job is located"
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
variable "lts_group_name" {
  description = "The name of the LTS group used to store the DRS job logs"
  type        = string
}

variable "lts_stream_name" {
  description = "The name of the LTS stream used to store the DRS job logs"
  type        = string
}

variable "drs_job_id" {
  description = "The ID of the existing DRS job"
  type        = string
}

variable "compare_policy_period" {
  description = "The comparison period of the compare policy, e.g. * * 1,3,5 for weekly comparison"
  type        = string
}

variable "compare_policy_begin_time" {
  description = "The start time when the comparison policy takes effect, UTC time in HH:mm:ss format"
  type        = string
}

variable "compare_policy_end_time" {
  description = "The end time when the comparison policy takes effect, UTC time in HH:mm:ss format"
  type        = string
}

variable "compare_policy_compare_type" {
  description = "The list of comparison types, valid values are object_comparison, lines and account"
  type        = list(string)
  default     = ["lines"]
}

variable "compare_policy_compare_policy" {
  description = "The comparison policy, valid values are normal and manyToOne"
  type        = string
  default     = "normal"
}

variable "compare_policy_interval_hour" {
  description = "The comparison interval in hours, required for hourly comparison"
  type        = number
  default     = null
}
