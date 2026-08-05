# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DSC mask algorithm is located"
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
variable "mask_algorithm_name" {
  description = "The name of the mask algorithm"
  type        = string
}

variable "mask_algorithm_prefix_length" {
  description = "The number of characters to retain at the beginning of the data"
  type        = number
  default     = 6
}

variable "mask_algorithm_suffix_length" {
  description = "The number of characters to retain at the end of the data"
  type        = number
  default     = 4
}

variable "mask_algorithm_replacement" {
  description = "The character used to mask the data"
  type        = string
  default     = "*"
}
