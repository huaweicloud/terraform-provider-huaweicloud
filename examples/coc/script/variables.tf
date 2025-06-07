variable "script_name" {
  description = "The name of the COC script"
  type        = string
}

variable "script_description" {
  description = "The description of the COC script"
  type        = string
}

variable "script_parameters" {
  description = "The parameter list of the COC script"
  type = list(object({
    name        = string
    value       = string
    description = string
    sensitive   = optional(bool)
  }))
}

