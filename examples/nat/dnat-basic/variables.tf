variable "ecs_name" {
  type    = string
  default = "ECS_Example"
}

variable "user_password" {
  type    = string
  default = "Base!!52259"
}

variable "vpc_name" {
  type    = string
  default = "vpc_Example"
}

variable "vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

variable "ims_name" {
  type    = string
  default = "CentOS 7.3 64bit"
}

variable "bandwidth_name" {
  type    = string
  default = "Bandwidth_Example"
}

variable "subnet_name" {
  type    = string
  default = "subnet_Example"
}

variable "subnet_cidr" {
  type    = string
  default = "192.168.64.0/18"
}

variable "net_gateway_name" {
  type    = string
  default = "net_gateway_Example"
}

variable "subnet_gateway_ip" {
  type    = string
  default = "192.168.64.1"
}

variable "ecs_ipaddress" {
  type    = string
  default = "192.168.64.15"
}

variable "secgroup_name" {
  type    = string
  default = "secgroup_Example"
}

variable "example_security_group" {
  type = list(object({
    direction         = string
    ethertype         = string
    protocol          = string
    port_range_min    = number
    port_range_max    = number
    remote_ip_prefix  = string
  }))
  default = [
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=null, port_range_max=null, remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv4", direction="ingress", protocol="icmp", port_range_min=null, port_range_max=null, remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv6", direction="ingress", protocol="tcp",  port_range_min=null, port_range_max=null, remote_ip_prefix="::/0"},
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=20,   port_range_max=21,   remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=22,   port_range_max=22,   remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=80,   port_range_max=80,   remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=443,  port_range_max=443,  remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=3389, port_range_max=3389, remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv4", direction="ingress", protocol="tcp",  port_range_min=8080, port_range_max=8080, remote_ip_prefix="0.0.0.0/0"},
    
    {ethertype="IPv4", direction="egress",  protocol="tcp",  port_range_min=null, port_range_max=null, remote_ip_prefix="0.0.0.0/0"},
    {ethertype="IPv6", direction="egress",  protocol="tcp",  port_range_min=null, port_range_max=null, remote_ip_prefix="::/0"}
  ]
}

variable "example_dnat_rule" {
  type = list(object({
    private_ip            = string
    internal_service_port = number
    protocol              = string
    external_service_port = number
  }))
  default = [
    {private_ip="192.168.64.15", internal_service_port=8080,  protocol="tcp", external_service_port=8080},
    {private_ip="192.168.64.15", internal_service_port=22,    protocol="tcp", external_service_port=8022},
    ]
}

variable "evs_name" {
  type    = string
  default = "volume_Example"
}