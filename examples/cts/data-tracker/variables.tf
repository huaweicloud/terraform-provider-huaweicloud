# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CTS service is located"
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

# Variable definitions for OBS resources
variable "source_bucket_name" {
  description = "The name of the OBS bucket for storing trace files"
  type        = string
}

variable "transfer_bucket_name" {
  description = "The name of the OBS bucket for transferring trace files"
  type        = string
}

# Variable definitions for CTS resources
variable "tracker_name" {
  description = "The name of the system tracker"
  type        = string
}

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
