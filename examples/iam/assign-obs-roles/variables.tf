variable "group_name" {
  type        = string
  description = "The name of the IAM identity group."
}

variable "account_domain_id" {
  type        = string
  description = "Account (domain) ID of the authorized account."
}

variable "user_name" {
  type        = string
  description = "user (IAM) name of the authorized account."
}
