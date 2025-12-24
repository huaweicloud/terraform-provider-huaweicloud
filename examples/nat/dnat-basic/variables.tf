variable "ecs_name" {
  type    = string
  default = "ECS_Example"
}

variable "ims_name" {
  type    = string
  default = "CentOS 7.3 64bit"
}

variable "vpc_name" {
  type    = string
  default = "vpc_Example"
}

variable "subnet_name" {
  type    = string
  default = "subnet_Example"
}

variable "secgroup_name" {
  type    = string
  default = "secgroup_Example"
}

variable "bandwidth_name" {
  type    = string
  default = "Bandwidth_Example"
}

variable "net_gateway_name" {
  type    = string
  default = "net_gateway_Example"
}

variable "vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

variable "subnet_cidr" {
  type    = string
  default = "192.168.64.0/18"
}

variable "subnet_gateway_ip" {
  type    = string
  default = "192.168.64.1"
}

variable "security_group_rule" {
  type    = list(object({
    direction        = string
    ethertype        = string
    protocol         = string
    port_range_min   = number
    port_range_max   = number
    remote_ip_prefix = string
  }))
  default = [
    { direction = "ingress", ethertype = "IPv4", protocol = "tcp", port_range_min = 80, port_range_max = 80, remote_ip_prefix = "0.0.0.0/0" },
  ]
}

variable "example_dnat_rule" {
  type    = list(object({
    private_ip            = string
    internal_service_port = number
    protocol              = string
    external_service_port = number
  }))
  default = [
    { private_ip = "192.168.64.15", internal_service_port = 80, protocol = "tcp", external_service_port = 8080 },
    { private_ip = "192.168.64.15", internal_service_port = 22, protocol = "tcp", external_service_port = 8022 },
  ]
}

variable "ecs_ssh_port" {
  type    = number
  default = 8022
}
