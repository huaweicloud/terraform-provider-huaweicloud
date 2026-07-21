# Variable definitions for authentication
variable "region_name" {
  description = "The region where the GaussDB instance is located"
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
variable "instance_id" {
  description = "The ID of the GaussDB instance"
  type        = string
  default     = ""
}

variable "node_ids" {
  description = "The list of node IDs to query"
  type        = list(string)
  default     = []
}

variable "start_time" {
  description = "The start time for the query, in 13-digit UNIX timestamp format (milliseconds, UTC)"
  type        = number
  default     = 0
}

variable "end_time" {
  description = "The end time for the query, in 13-digit UNIX timestamp format (milliseconds, UTC)"
  type        = number
  default     = 0
}

variable "support_system" {
  description = "Whether to display system users"
  type        = bool
  default     = false
}

variable "multi_queries" {
  description = "The list of field aggregation query conditions"
  type        = list(object({
    name      = string
    condition = string
    values    = list(string)
    is_fuzzy  = optional(string)
  }))
  default     = []
}
