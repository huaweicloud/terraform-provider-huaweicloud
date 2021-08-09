variable "vpc_name" {
  description = "The name of the Huaweicloud VPC"
  default     = "tf_vpc_demo"
}

variable "subnet_name" {
  description = "The name of the Huaweicloud Subnet"
  default     = "tf_subnet_demo"
}

variable "security_group_name" {
  description = "The name of the Huaweicloud Security Group"
  default     = "tf_secgroup_demo"
}

# If you want to use other operation systems, please update the nginx install commands of the provisioner.
variable "ecs_image" {
  description = "The name of the Huaweicloud Image"
  default     = "CentOS 7.3 64bit"
}

variable "ecs_name" {
  description = "The instance name of the Huaweicloud ECS"
  default     = "tf_ecs_demo"
}

variable "bandwidth_name" {
  description = "The bandwidth name of the Huaweicloud EIP"
  default     = "tf_bandwidth_demo"
}

variable "function_name" {
  description = "The function name of the Huaweicloud FunctionGraph"
  default     = "tf_function_demo"
}

variable "apig_instance_name" {
  description = "The instance name of the Huaweicloud Dedicated APIG"
  default     = "tf_apig_instance_demo"
}

variable "apig_auth_name" {
  description = "The custom authorizer name of the Huaweicloud Dedicated APIG"
  default     = "tf_apig_auth_demo"
}

variable "apig_response_name" {
  description = "The response name of the Huaweicloud Dedicated APIG"
  default     = "tf_apig_response_demo"
}

variable "apig_group_name" {
  description = "The group name of the Huaweicloud Dedicated APIG"
  default     = "tf_apig_group_demo"
}

variable "apig_channel_name" {
  description = "The channel name of the Huaweicloud Dedicated APIG"
  default     = "tf_vpc_channel_demo"
}

variable "apig_api_name" {
  description = "The API name of the Huaweicloud Dedicated APIG"
  default     = "tf_apig_api_demo"
}
