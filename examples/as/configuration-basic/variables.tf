variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "key_pair_name" {
  description = "The name of the key pair"
  type        = string
}

variable "public_key" {
  description = "The public key for the key pair"
  type        = string
  sensitive   = true
}

variable "configuration_name" {
  description = "The name of the AS configuration"
  type        = string
}
