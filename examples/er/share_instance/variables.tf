variable "region_name" {
  default = "cn-north-4"
}

variable "owner_account_id" {
  type        = string
  description = "The account ID of the ER instance sharer"

  default   = ""
  sensitive = true
}

variable "owner_ak" {
  type        = string
  description = "The AK of the ER instance sharer"

  default   = ""
  sensitive = true
}

variable "owner_sk" {
  type        = string
  description = "The SK of the ER instance sharer"

  default   = ""
  sensitive = true
}

variable "principal_account_id" {
  type        = string
  description = "The account ID of the ER instance accepter"

  default   = ""
  sensitive = true
}


variable "principal_ak" {
  type        = string
  description = "The AK of the ER instance accepter"

  default   = ""
  sensitive = true
}

variable "principal_sk" {
  type        = string
  description = "The sk of the ER instance accepter"

  default   = ""
  sensitive = true
}

variable "er_instance_name" {
  type        = string
  description = "The ID of the shared ER instance"

  default = "share_attachment"
}

variable "resource_share_name" {
  type        = string
  description = "The ID of the shared ER instance"

  default = "resource-share-er"
}

variable "vpc_name" {
  type        = string
  description = "The ID of the VPC of the VPC attachment"

  default = "er_attachment_vpc_test"
}

variable "subnet_name" {
  type        = string
  description = "The ID of the subnet of the VPC attachment"

  default = "er_attachment_subnet_test"
}

variable "attachment_name" {
  type        = string
  description = "The ID of the VPC attachment"

  default = "shared_attachment"
}
