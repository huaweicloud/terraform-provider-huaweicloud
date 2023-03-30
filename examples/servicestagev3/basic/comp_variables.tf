variable "comp_name" {
  default = "comp_terraform_1"
}

variable "comp_description" {
  default = "comp_description"
}

variable "comp_version" {
  default = "1.1.1"
}

variable "comp_environment_id" {
  type = string
}

variable "comp_application_id" {
  type = string
}

variable "comp_replica" {
  type = number
}

variable "comp_runtime_stack" {
  type    = object({ name = string, version = string, type = string, deploy_mode = string })
  default = {
    name        = "OpenJDK8"
    version     = "1.1.1"
    type        = "Java"
    deploy_mode = "virtualmachine"
  }
}
