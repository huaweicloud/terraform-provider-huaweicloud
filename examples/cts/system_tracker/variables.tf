# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CTS service is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for OBS resources
variable "bucket_name" {
  description = "The name of the OBS bucket for storing trace files"
  type        = string
}

# Variable definitions for CTS resources
variable "tracker_enabled" {
  description = "Whether to enable the system tracker"
  type        = bool
  default     = true
}

variable "tracker_tags" {
  description = "The tags of the system tracker"
  type        = map(string)
  default     = {}
}

variable "is_system_tracker_delete" {
  description = "Whether to delete the system tracker when the tracker resource is deleted"
  type        = bool
  default     = true
}

variable "trace_object_prefix" {
  description = "The prefix of the trace object in the OBS bucket"
  type        = string
}

variable "trace_file_compression_type" {
  description = "The compression type of the trace file"
  type        = string
  default     = "gzip"
}

variable "is_lts_enabled" {
  description = "Whether to enable the trace analysis for LTS service"
  type        = bool
  default     = true
}
