variable "elb_loadbalancer_name" {
  description = "The name of the Huaweicloud ELB loadbalancer"
  default     = "tf_elb_loadbalancer_demo"
}

variable "elb_listener_name" {
  description = "The name of the Huaweicloud ELB listener"
  default     = "tf_elb_listener_demo"
}

variable "vpc_name" {
  description = "The name of the Huaweicloud VPC"
  default     = "tf_vpc_demo"
}

variable "vpc_cidr" {
  description = "The CIDR of the Huaweicloud VPC"
  default     = "172.16.0.0/16"
}

variable "subnet_name" {
  description = "The name of the Huaweicloud VPC subnet"
  default     = "tf_subnet_demo"
}

variable "subnet_cidr" {
  description = "The CIDR of the Huaweicloud VPC subnet"
  default     = "172.16.10.0/24"
}

variable "subnet_gateway" {
  description = "The gateway IP address of the Huaweicloud VPC subnet"
  default     = "172.16.10.1"
}

variable "security_group_name" {
  description = "The name of the Huaweicloud Security group"
  default     = "tf_security_group_demo"
}

variable "ecs_instance_name" {
  description = "The name of the Huaweicloud ECS instance"
  default     = "tf_ecs_instance_demo"
}

variable "network_acl_name" {
  description = "The name of the Huaweicloud Network ACL"
  default     = "tf_network_acl_demo"
}

variable "network_acl_rule_name" {
  description = "The name of the Huaweicloud Network ACL rule"
  default     = "tf_network_acl_rule_demo"
}
