# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CBH instance is located"
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
variable "cluster_id" {
  description = "The ID of the CCE cluster"
  type        = string
  default     = ""

  validation {
    condition     = var.cluster_id != "" || var.cluster_name != ""
    error_message = "One of cluster_id or cluster_name is required"
  }
}

variable "cluster_name" {
  description = "The name of the CCE cluster"
  type        = string
  default     = ""
}

variable "addon_template_name" {
  description = "The name of the CCE addon template"
  type        = string
  default     = "coredns"
}

variable "addon_version" {
  description = "The version of the CCE addon template"
  type        = string
}
