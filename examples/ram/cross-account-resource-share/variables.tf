# Variable definitions for authentication
variable "region_name" {
  description = "The region where the RAM resource share is located"
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

# Resource share variables
variable "resource_share_name" {
  description = "The name of the resource share"
  type        = string
}

variable "description" {
  description = "The description of the resource share"
  type        = string
  default     = ""
}

variable "principals" {
  description = "The list of one or more principals (account IDs or organization IDs) to share resources with"
  type        = list(string)
}

variable "resource_urns" {
  description = "The list of URNs of one or more resources to be shared. If not specified, URNs will be automatically generated from created resources (VPC, subnet, security group)"
  type        = set(string)
  default     = []
}

variable "permission_ids" {
  description = "The list of RAM permissions associated with the resource share"
  type        = list(string)
  default     = []
}

variable "allow_external_principals" {
  description = "Whether resources can be shared with any accounts outside the organization"
  type        = bool
  default     = false
}
