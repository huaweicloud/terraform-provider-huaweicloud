# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ER resources and VPC resources are located"
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

variable "principal_access_key" {
  description = "The access key of the principal IAM user"
  type        = string
  sensitive   = true
}

variable "principal_secret_key" {
  description = "The secret key of the principal IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for authentication
variable "instance_name" {
  description = "The name of the ER instance"
  type        = string
}

variable "instance_asn" {
  description = "The ASN of the ER instance"
  type        = number
  default     = 64512
}

variable "instance_description" {
  description = "The description of the ER instance"
  type        = string
  default     = "The ER instance to share with other accounts"
}

variable "instance_enable_default_propagation" {
  description = "Whether to enable the default propagation"
  type        = bool
  default     = true
}

variable "instance_enable_default_association" {
  description = "Whether to enable the default association"
  type        = bool
  default     = true
}

variable "instance_auto_accept_shared_attachments" {
  description = "Whether to automatically accept the shared attachments"
  type        = bool
  default     = false
}

variable "resource_share_name" {
  description = "The ID of the shared ER instance"
  type        = string

  default = "resource-share-er"
}

variable "principal_account_id" {
  description = "The account ID of the ER instance accepter"
  type        = string
}

variable "owner_account_id" {
  description = "The account ID of the ER instance sharer"
  type        = string
}

variable "principal_vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "principal_vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "principal_subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "principal_subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "principal_subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "attachment_name" {
  description = "The name of the attachment"
  type        = string
}
