variable "security_group_id" {
  description = "The id of the security group in which to create the loadbalancer."
}

variable "vpc_id" {
  description = "Specifies the VPC ID."
}

variable "vip_subnet_id" {
  description = "Specifies the ID of the subnet to which a VIP is assigned. Only an internal network is supported."
}

variable "vip_address" {
  description = "Specifies the IP address used by ELB for providing services. When type is set to External, the parameter value is the elastic IP address. When type is set to Internal, the parameter value is the private network IP address."
}

variable "tenantid" {
  description = "Specifies the tenant ID. This parameter is mandatory when type is set to Internal."
}

variable "available_zone" {
  description = "This field is valid when type is set to Internal. If type is set to Internal and an AZ is specified, the specified AZ must support private network load balancers. Otherwise, an error message is returned."
  default     = "The first az of the az list"
}

variable "vm_private_address" {
  description = "Specifies the private IP address of the backend ECS."
}

variable "vm_id" {
  description = "Specifies the backend ECS ID."
}
