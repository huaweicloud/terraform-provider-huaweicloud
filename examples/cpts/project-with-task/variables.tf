# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CPTS service is located"
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

# Variable definitions for CPTS resources
variable "cpts_project_name" {
  description = "The name of the CPTS project"
  type        = string
}

variable "cpts_project_description" {
  description = "The description of the CPTS project"
  type        = string
  default     = ""
}

variable "cpts_task_name" {
  description = "The name of the CPTS test task"
  type        = string
}

variable "cpts_task_benchmark_concurrency" {
  description = "The benchmark concurrency of the CPTS test task"
  type        = number
  default     = 200
}
