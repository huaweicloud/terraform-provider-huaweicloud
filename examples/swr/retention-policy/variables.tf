# Variable definitions for authentication
variable "region_name" {
  description = "The region where the SWR service is located"
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

variable "organization_name" {
  description = "The organization name"
  type        = string
}

variable "repository_name" {
  description = "The repository name"
  type        = string
}

variable "category" {
  description = "The category"
  type        = string
  default     = "linux"
}

variable "policy_type" {
  description = "The policy type"
  type        = string
}

variable "policy_number" {
  description = "The policy number"
  type        = number
}

variable "tag_selectors_configuration" {
  description = "The configuration of the tag selectors"
  type        = list(object({
    kind    = string
    pattern = number
  }))

  default = []
}
