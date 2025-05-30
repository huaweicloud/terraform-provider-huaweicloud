variable "script_name" {
  description = "The name of the COC script"
  default     = "tf_coc_script_name"
}

variable "script_execute_name" {
  description = "The COC script execute name"
  default     = "tf_script_execute_name"
}

variable "operation_type" {
  description = "The COC script order operation type"
  default     = "CANCEL_ORDER"
}
