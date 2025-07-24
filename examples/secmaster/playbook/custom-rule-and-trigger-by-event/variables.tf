# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SecMaster playbook is located"
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

# Variable definitions for resources
variable "workspace_id" {
  description = "The ID of the SecMaster workspace"
  type        = string
  default     = ""
}

variable "workspace_name" {
  description = "The name of the SecMaster workspace"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.workspace_id != "" || var.workspace_name != ""
    error_message = "At least one of workspace_id and workspace_name must be provided."
  }
}

variable "playbook_name" {
  description = "The name of the SecMaster playbook"
  type        = string
}

variable "rule_expression_type" {
  description = "The expression type of the playbook rule"
  type        = string
  default     = "custom"
}

variable "rule_conditions" {
  description = "The condition rule list of the playbook"
  type        = list(object({
    name   = string
    detail = string
    data   = list(string)
  }))

  validation {
    condition     = length(var.rule_conditions) >= 2
    error_message = "The length of rule_conditions must be greater than or equal to 2."
  }
}

variable "approval_content" {
  description = "The approval content for the playbook version"
  type        = string
  default     = "Approved for production use"
}
