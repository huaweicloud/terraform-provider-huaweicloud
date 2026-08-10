# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DIS stream is located"
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

# Variable definitions for huaweicloud_dis_stream
variable "stream_name" {
  description = "The name of the DIS stream"
  type        = string
}

variable "stream_partition_count" {
  description = "The number of partitions for the DIS stream"
  type        = number
}

variable "stream_type" {
  description = "The type of the DIS stream"
  type        = string
  default     = null
}

variable "stream_retention_period" {
  description = "The data retention period in hours"
  type        = number
  default     = 24
}

variable "stream_auto_scale_min_partition_count" {
  description = "The minimum number of partitions for auto scaling"
  type        = number
  default     = null
}

variable "stream_auto_scale_max_partition_count" {
  description = "The maximum number of partitions for auto scaling"
  type        = number
  default     = null
}

variable "stream_compression_format" {
  description = "The compression format of the data"
  type        = string
  default     = null
}

variable "stream_data_type" {
  description = "The type of the data"
  type        = string
  default     = null
}

variable "stream_csv_delimiter" {
  description = "The delimiter for CSV data"
  type        = string
  default     = null
}

variable "stream_data_schema" {
  description = "The schema of the data"
  type        = string
  default     = null
}

variable "stream_tags" {
  description = "The key/value pairs to associate with the DIS stream"
  type        = map(string)
  default     = {}
}
