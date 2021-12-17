
variable "is_template" {
   default = true
}

variable "volumetype" {
   default = "SAS"
}

variable "data_volume_type" {
   default = "SAS"
}

variable "vpc_cidr" {
   default = "192.168.0.0/16"
}

variable "publicip_type" {
   default = "5_bgp"
}

variable "nics_cidr" {
   default = "192.168.0.0/16"
}