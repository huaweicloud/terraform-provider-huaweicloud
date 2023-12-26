variable "vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

variable "subnetA_cidr" {
  type    = string
  default = "192.168.1.0/24"
}

variable "subnetA_gateway_ip" {
  type    = string
  default = "192.168.1.1"
}

variable "subnetB_cidr" {
  type    = string
  default = "192.168.2.0/24"
}

variable "subnetB_gateway_ip" {
  type    = string
  default = "192.168.2.1"
}

variable "peer_conn_name" {
  type    = string
  default = "peering_test"
}
