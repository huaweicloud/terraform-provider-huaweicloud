# Variable definitions for authentication
variable "region_name" {
  description = "The region where the FunctionGraph service is located"
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

# Variable definitions for FunctionGraph resources
variable "function_name" {
  description = "The name of the FunctionGraph function"
  type        = string
}

variable "function_memory_size" {
  description = "The memory size of the function in MB"
  type        = number
  default     = 128
}

variable "function_timeout" {
  description = "The timeout of the function in seconds"
  type        = number
  default     = 10
}

variable "function_runtime" {
  description = "The runtime of the function"
  type        = string
  default     = "Python2.7"
}

variable "function_code" {
  description = "The source code of the function"
  type        = string
  default     = <<EOT
import json

def handler(event, context):
    print("FunctionGraph timer trigger executed!")
    print("Event:", json.dumps(event))
    print("Context:", json.dumps(context))
    return {
        'statusCode': 200,
        'body': json.dumps('Hello, FunctionGraph!')
    }
  EOT
}

variable "function_description" {
  description = "The description of the function"
  type        = string
  default     = ""
}

variable "trigger_status" {
  description = "The status of the FunctionGraph timer trigger"
  type        = string
  default     = "ACTIVE"
}

variable "trigger_name" {
  description = "The name of the FunctionGraph timer trigger"
  type        = string
}

variable "trigger_schedule_type" {
  description = "The schedule type of the FunctionGraph timer trigger"
  type        = string
  default     = "Cron"
}

variable "trigger_sync_execution" {
  description = "Whether to execute the function synchronously"
  type        = bool
  default     = false
}

variable "trigger_user_event" {
  description = "The user event description for the FunctionGraph timer trigger"
  type        = string
  default     = ""
}

variable "trigger_schedule" {
  description = "The schedule expression for the FunctionGraph timer trigger"
  type        = string
  default     = "@every 1h30m"
}
