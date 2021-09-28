variable "keypair_name" {
  description = "Keypair name"
}

variable "private_key_path" {
  description = "The relative path of the private key"
}

variable "vpc_name" {
  description = "Name of the HuaweiCloud VPC"
}

variable "vpc_cidr" {
  description = "CIDR of the HuaweiCloud VPC"
}

variable "subnet_name" {
  description = "Name of the HuaweiCloud VPC Subnet"
}

variable "subnet_cidr" {
  description = "CIDR of the HuaweiCloud VPC subnet"
}

variable "security_group_name" {
  description = "Name of the HuaweiCloud Network security group"
}

variable "gateway_ip" {
  description = "Gateway IP of the HuaweiCloud VPC subnet"
}

variable "image_name" {
  description = "Name of the HuaweiCloud Image"
  default     = "Ubuntu 18.04 server 64bit"
}

variable "ecs_instance_name" {
  description = "Name of the HuaweiCloud ECS instance"
}

variable "bandwidth_name" {
  description = "Bandwidth name of the HuaweiCloud EIP"
}