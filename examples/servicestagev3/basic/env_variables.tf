variable "env_name" {
  default = "env_terraform_1"
}

variable "env_deploy_mode" {
  default = "virtualmachine"
}

variable "env_description" {
  default = "env_description"
}

variable "env_labels" {
  type    = list(object({ key = string, value = string }))
  default = [
    { key = "labels_key1", value = "value1" },
    { key = "labels_key2", value = "value2" },
  ]
}

variable "resources" {
  type    = list(object({ id = string, name = string, type = string }))
}
