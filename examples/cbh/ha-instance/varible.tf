variable name {
  type        = string
  default     = "Cbh_HA_demo"
  description = "The name of the CBH(HA) instance."
}

variable flavor_id {
  type        = string
  default     = "cbh.basic.50"
  description = "The product ID of the CBH(HA) server."
}

variable password {
  type        = string
  default     = "Cbh@Huawei123"
  description = "The password for logging in the management console."
  sensitive   = true
}

variable charging_mode {
  type        = string
  default     = "prePaid"
  description = "The charging mode of the CBH(HA) instance."
}

variable period_unit {
  type        = string
  default     = "month"
  description = "The charging period unit of the CBH(HA) instance."
}

variable period {
  type        = number
  default     = 1
  description = "The charging period of the CBH(HA) instance."
}
