# Variable definitions for authentication
variable "region_name" {
  description = "The region where the workspace is located"
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
variable "workspace_id" {
  description = "The ID of the workspace"
  type        = string
  default     = ""
  nullable    = false
}

variable "workspace_name" {
  description = "The name of the workspace"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.workspace_id != "" || var.workspace_name != ""
    error_message = "At least one of workspace_id and workspace_name must be provided."
  }
}

variable "workflow_id" {
  description = "The ID of the workflow"
  type        = string
  default     = ""
  nullable    = false
}

variable "workflow_name" {
  description = "The name of the workflow"
  type        = string
}

variable "workflow_version_taskflow" {
  description = "The Base64 encoded of the workflow topology diagram"
  type        = string
}

variable "workflow_version_taskconfig" {
  description = "The parameters configuration of the workflow topology diagram"
  type        = string
}

variable "workflow_version_taskflow_type" {
  description = "The taskflow type of the workflow"
  type        = string
  default     = "JSON"
}

variable "workflow_version_aop_type" {
  description = "The aop type of the workflow"
  type        = string
  default     = "NORMAL"
}

variable "workflow_version_description" {
  description = "The description of the workflow version"
  type        = string
  default     = ""
}
