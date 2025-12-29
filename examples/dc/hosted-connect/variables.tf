# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
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

# Variable definitions for hosted connect
variable "hosted_connect_name" {
  description = "The name of the hosted connect"
  type        = string
}

variable "hosted_connect_description" {
  description = "The description of the hosted connect"
  type        = string
  default     = "Created by Terraform"
}

variable "bandwidth" {
  description = "The bandwidth size of the hosted connect in Mbit/s"
  type        = number
}

variable "hosting_id" {
  description = "The ID of the operations connection on which the hosted connect is created"
  type        = string
}

variable "vlan" {
  description = "The VLAN allocated to the hosted connect"
  type        = number
}

variable "resource_tenant_id" {
  description = "The tenant ID for whom a hosted connect is to be created"
  type        = string
}

variable "peer_location" {
  description = "The location of the on-premises facility at the other end of the connection"
  type        = string
  default     = ""
}
