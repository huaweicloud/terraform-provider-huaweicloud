# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RAM service is located"
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
variable "resource_share_ids" {
  description = "List of resource share IDs to query invitations for"
  type        = list(string)
  default     = []
}

variable "action" {
  description = "The action to perform on invitations"
  type        = string
  default     = "reject"

  validation {
    condition     = contains(["accept", "reject"], var.action)
    error_message = "The action must be either 'accept' or 'reject'."
  }
}
