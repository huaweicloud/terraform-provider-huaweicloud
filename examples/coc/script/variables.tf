variable "coc_script_name" {
  description = "The name of the script"
  type        = string
}

variable "coc_script_description" {
  description = "The description of the script"
  type        = string
}

variable "coc_script_risk_level" {
  description = "The description of the script"
  type        = string
}

variable "coc_script_version" {
  description = "The description of the script"
  type        = string
}

variable "coc_script_type" {
  description = "The risk level of the script"
  type        = string
}

variable "coc_script_content" {
  description = "The content of the script"
  type        = string
}

variable "coc_script_parameters" {
  description = "The parameter list of the script"
  type = list(object({
    name        = string
    value       = string
    description = string
    sensitive   = optional(bool)
  }))
}

