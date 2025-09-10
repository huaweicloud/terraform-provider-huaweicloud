# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CTS service is located"
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

# Variable definitions for SMN resources
variable "topic_name" {
  description = "The name of the SMN topic"
  type        = string
}

# Variable definitions for CTS resources
variable "notification_name" {
  description = "The name of the CTS notification"
  type        = string
}

variable "notification_operation_type" {
  description = "The type of operation"
  type        = string
  default     = "customized"
}

variable "notification_agency_name" {
  description = "The name of the IAM agency which allows CTS to access the SMN resources"
  type        = string
}

variable "notification_filter" {
  description = "The filter of the notification"
  type        = list(object({
    condition = string
    rule      = list(string)
  }))

  default  = []
  nullable = false
}

variable "notification_operations" {
  description = "The operations of the notification"
  type        = list(object({
    service     = string
    resource    = string
    trace_names = list(string)
  }))

  default  = []
  nullable = false
}

variable "notification_operation_users" {
  description = "The operation users of the notification"
  type        = list(object({
    group = string
    users = list(string)
  }))

  default  = []
  nullable = false
}
