variable "source_obs_name" {
  default     = "oms-test-source"
  description = "The name of source obs"
}

variable "source_dest_name" {
  default     = "oms-test-dest"
  description = "The name of dest obs"
}

variable "source_region" {
  default     = "cn-north-4"
  description = "The huaweicloud obs source region"
}

variable "dest_region" {
  default     = "cn-south-1"
  description = "The huaweicloud obs destination region"
}

variable "source_ak" {
  description = "The access_key of source huaweicloud obs"
}

variable "source_sk" {
  description = "The secret_key of source huaweicloud obs"
}

variable "dest_ak" {
  description = "The access_key of destination huaweicloud obs"
}

variable "dest_sk" {
  description = "The secret_key of destination huaweicloud obs"
}