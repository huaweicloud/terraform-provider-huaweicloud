# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SecMaster resources will be created"
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

# Variable definitions for configuration dictionary
variable "dict_id" {
  description = "The dictionary ID"
  type        = string
}

variable "dict_key" {
  description = "The dictionary key"
  type        = string
}

variable "dict_code" {
  description = "The dictionary code"
  type        = string
}

variable "dict_val" {
  description = "The dictionary value"
  type        = string
}

variable "language" {
  description = "The language environment. Valid values: zh, en"
  type        = string
  default     = "zh"
}

variable "dict_version" {
  description = "The version number"
  type        = string
  default     = "1.0.0"
}

variable "dict_pkey" {
  description = "The parent key of the dictionary"
  type        = string
  default     = ""
}

variable "dict_pcode" {
  description = "The parent code of the dictionary"
  type        = string
  default     = ""
}

variable "scope" {
  description = "The domain to which the dictionary belongs"
  type        = string
  default     = "ALERT"
}

variable "description" {
  description = "The description of the dictionary"
  type        = string
  default     = ""
}

variable "extend_field" {
  description = "The extension field of the dictionary"
  type        = map(string)
  default     = {}
}
