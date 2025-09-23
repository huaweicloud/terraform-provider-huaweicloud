# Variable definitions for authentication
variable "region_name" {
  description = "The region where the VPC is located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
}

# Variable definitions for resources/data sources
variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "172.16.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = "172.16.10.0/24"
}

variable "subnet_gateway" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = "172.16.10.1"
}

variable "channel_name" {
  description = "The name of the custom event channel"
  type        = string
}

variable "sources_name" {
  description = "The name of the event source"
  type        = string
  default     = "HC.VPC"
}

variable "sources_provider_type" {
  description = "The provider type of the event source"
  type        = string
  default     = "OFFICIAL"
}

variable "source_op" {
  description = "The operation of the source"
  type        = string
  default     = "StringIn"
}

variable "type_op" {
  description = "The operation of the type"
  type        = string
  default     = "StringIn"
}

variable "subscription_source_values" {
  description = "The event types to be subscribed from VPC service"
  type        = list(string)
  default     = [
    "VPC:CloudTrace:ApiCall",
    "VPC:CloudTrace:ConsoleAction",
    "VPC:CloudTrace:SystemAction"
  ]
}

variable "targets_name" {
  description = "The name of the event target"
  type        = string
  default     = "HC.EG"
}

variable "targets_provider_type" {
  description = "The type of the event target"
  type        = string
  default     = "OFFICIAL"
}

variable "detail_name" {
  description = "The name(key) of the target detail configuration"
  type        = string
  default     = "eg_detail"
}

variable "transform" {
  description = "The transform configuration of the event target, in JSON format"
  type        = map(string)
  default     = {
    "type" : "ORIGINAL",
  }
}

variable "agency_name" {
  description = "The name of the agency"
  type        = string
  default     = "EG_AGENCY"
}

variable "target_project_id" {
  description = "The ID of the target project"
  type        = string
}

variable "target_region_name" {
  description = "The name of the target region"
  type        = string
  default     = ""
}
