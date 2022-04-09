variable "iden_user_name" {
  type    = string
  default = "user_A"
}

variable "iden_group" {
  description = "List of identity group."
  type = list(object({
    name        = string
    description = string
  }))
  default = [
    { name = "group_A", description = "This is a identity group." },
    { name = "group_B", description = "This is a identity group." }
  ]
}

variable "domain_id" {
  type        = string
  default     = "your_domain_id"
  description = "This is the domain id."
}
