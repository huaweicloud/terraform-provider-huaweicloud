# Variable definitions for authentication
variable "region_name" {
  description = "The region where the EIP resources are located"
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
variable "enterprise_project_id" {
  description = "The ID of the enterprise project"
  type        = string
  default     = ""
  nullable    = false
}

variable "bandwidth_name" {
  description = "The name of the shared bandwidth"
  type        = string
}

variable "bandwidth_size" {
  description = "The size of the shared bandwidth in Mbit/s"
  type        = number
  default     = 5
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the shared bandwidth"
  type        = string
  default     = "bandwidth"
}

variable "bandwidth_type" {
  description = "The type of the bandwidth"
  type        = string
  default     = "share"
}

variable "bandwidth_public_border_group" {
  description = "The border group of the public IP"
  type        = string
  default     = "center"
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "eip_bandwidth_name" {
  description = "The name of the dedicated EIP bandwidth"
  type        = string
}

variable "eip_bandwidth_size" {
  description = "The size of the dedicated EIP bandwidth in Mbit/s"
  type        = number
  default     = 5
}

variable "eip_bandwidth_charge_mode" {
  description = "The charge mode of the dedicated EIP bandwidth"
  type        = string
  default     = "traffic"
}

variable "eip_description" {
  description = "The description of the EIP"
  type        = string
  default     = ""
}

variable "eip_tags" {
  description = "The tags of the EIP"
  type        = map(string)
  default     = {}
}
