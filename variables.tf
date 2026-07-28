variable "access_key" {
  type  = string
  description = "HuaweiCloud API access key"
  sensitive = true
}

variable "secret_key"{
  type  = string
  description = "HuaweiCloud API secret key"
  sensitive = true
}

variable "region_name"{
  type  = string
  description = "HuaweiCloud region ()"
  default =  "cn-north-4"
}

variable  "security_token"{
  type  = string
  description = ""
  sensitive = true
}