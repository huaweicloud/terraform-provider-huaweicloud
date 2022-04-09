variable "vpc_name" {
  default = "vpc-basic"
}

variable "vpc_cidr" {
  default = "172.16.0.0/16"
}

variable "subnet_name" {
  default = "subent-basic"
}

variable "subnet_cidr" {
  default = "172.16.10.0/24"
}

variable "subnet_gateway" {
  default = "172.16.10.1"
}

variable "primary_dns" {
  default = "100.125.1.250"
}

variable "secondary_dns" {
  default = "100.125.21.250"
}

variable "bandwidth_name" {
  default = "mybandwidth"
}

variable "key_pair_name" {
  default = "mykey_pair"
}

variable "cce_cluster_name" {
  default = "mycce"
}

variable "cce_cluster_flavor" {
  default = "cce.s1.small"
}

variable "node_name" {
  default = "mynode"
}

variable "node_flavor" {
  default = "t6.large.2"
}

variable "root_volume_size" {
  default = 40
}

variable "root_volume_type" {
  default = "SAS"
}

variable "data_volume_size" {
  default = 100
}

variable "data_volume_type" {
  default = "SAS"
}

variable "ecs_flavor" {
  default = "sn3.large.2"
}

variable "ecs_name" {
  default = "myecs"
}

variable "os" {
  default = "EulerOS 2.5"
}

variable "image_name" {
  default = "EulerOS 2.5 64bit"
}
