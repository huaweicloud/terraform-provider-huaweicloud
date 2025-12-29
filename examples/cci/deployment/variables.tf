variable "region_name" {
  description = "The region where CCI deployment is located"
  type        = string
}

variable "access_key" {
  description = "The access key of IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of IAM user"
  type        = string
  sensitive   = true
}

variable "namespace_name" {
  description = "The name of CCI namespace"
  type        = string
}

variable "deployment_name" {
  description = "The name of CCI deployment"
  type        = string
}

variable "instance_type" {
  description = "The instance type of CCI pod"
  type        = string
  default     = "general-computing"
}

variable "container_name" {
  description = "The name of container"
  type        = string
  default     = "c1"
}

variable "container_image" {
  description = "The image of container"
  type        = string
  default     = "alpine:latest"
}

variable "cpu_limit" {
  description = "The CPU limit of container"
  type        = string
  default     = "1"
}

variable "memory_limit" {
  description = "The memory limit of container"
  type        = string
  default     = "2G"
}

variable "image_pull_secret_name" {
  description = "The name of image pull secret"
  type        = string
  default     = "imagepull-secret"
}
