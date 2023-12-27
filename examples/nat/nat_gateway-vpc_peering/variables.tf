variable "vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

variable "subnetA_gateway_ip" {
  type    = string
  default = "192.168.1.1"
}

variable "subnetB_gateway_ip" {
  type    = string
  default = "192.168.2.1"
}

variable "gateway_name" {
  type    = string
  default = "nat_gateway_test"
}

variable "secgroup_name" {
  type    = string
  default = "secgroup_test"
}

variable "remote_ip_prefix" {
  type = string
}

variable "image_name" {
  type    = string
  default = "Ubuntu 22.04 server 64bit"
}

variable "admin_pass" {
  type = string
}

variable "internal_service_port" {
  type    = string
  default = "22"
}

variable "external_service_port_A" {
  type    = string
  default = "8023"
}

variable "external_service_port_B" {
  type    = string
  default = "8823"
}

variable "peer_conn_name" {
  type    = string
  default = "peering_test"
}
