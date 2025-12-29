# Provider Configuration
variable "region_name" {
  description = "The name of the region to deploy resources"
  type        = string
}

variable "access_key" {
  description = "The access key of HuaweiCloud"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of HuaweiCloud"
  type        = string
  sensitive   = true
}

# CES dashboard Configuration
variable "dashboard_name" {
  description = "The name of the alarm dashboard"
  type        = string
}

variable "dashboard_row_widget_num" {
  description = "The monitoring view display mode"
  type        = number
}

variable "dashboard_extend_info" {
  description = "The information about the extension"
  type        = list(object({
    filter                  = string
    period                  = string
    display_time            = number
    refresh_time            = number
    from                    = number
    to                      = number
    screen_color            = string
    enable_screen_auto_play = bool
    time_interval           = number
    enable_legend           = bool
    full_screen_widget_num  = number
  }))
  default     = []
}
