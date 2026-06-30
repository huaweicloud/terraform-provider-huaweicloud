# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DataArts Studio instance is located"
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

# Variable definitions for DataArts Studio instance and workspace
variable "workspace_id" {
  description = "The ID of the workspace"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_id" {
  description = "The ID of the DataArts Studio instance to which the workspace belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "workspace_name" {
  description = "The name of the workspace used to filter results"
  type        = string
  default     = ""
  nullable    = false
}

# Variable definitions for data recognition rule group
variable "category_ids" {
  description = "The ID list of data categories used by data recognition rules"
  type        = list(string)
  default     = []
  nullable    = false
}

variable "data_recognition_rule_count" {
  description = "The number of data recognition rules to create and include in the rule group"
  type        = number
  default     = 1
}

variable "rule_group_name" {
  description = "The name of the data recognition rule group"
  type        = string
  default     = ""
  nullable    = false
}

variable "rule_group_description" {
  description = "The description of the data recognition rule group"
  type        = string
  default     = ""
}
