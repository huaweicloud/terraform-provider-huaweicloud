# Variable definitions for authentication
variable "region_name" {
  description = "The region where the resources are located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for RabbitMQ instance
variable "instance_flavor_type" {
  description = "The flavor type of the RabbitMQ instance"
  type        = string
  default     = "cluster"
}

variable "instance_storage_spec_code" {
  description = "The storage specification code of the RabbitMQ instance"
  type        = string
  default     = "dms.physical.storage.ultra.v2"
}

variable "ecs_image_name" {
  description = "The image name of the ECS instances (Ubuntu)"
  type        = string
  default     = "Ubuntu 20.04 server 64bit"
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = "192.168.0.0/24"
  nullable    = true
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = "192.168.0.1"
  nullable    = true
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "security_group_rule_configurations" {
  description = "The list of security group rule configurations."

  type = list(object({
    direction        = string
    ethertype        = string
    protocol         = string
    port_range_min   = number
    port_range_max   = number
    remote_ip_prefix = string
    description      = string
  }))

  default = [
    {
      direction        = "ingress"
      ethertype        = "IPv4"
      protocol         = "tcp"
      port_range_min   = 5672
      port_range_max   = 5672
      remote_ip_prefix = "192.168.0.0/16"
      description      = "Allow ECS instances to access RabbitMQ"
    },
    {
      direction        = "ingress"
      ethertype        = "IPv4"
      protocol         = "tcp"
      port_range_min   = 22
      port_range_max   = 22
      remote_ip_prefix = "192.168.0.0/16"
      description      = "Allow SSH access"
    }
  ]
}

variable "instance_name" {
  description = "The name of the RabbitMQ instance"
  type        = string
}

variable "instance_engine_version" {
  description = "The engine version of the RabbitMQ instance"
  type        = string
  default     = "3.8.35"
}

variable "instance_broker_num" {
  description = "The number of brokers of the RabbitMQ instance"
  type        = number
  default     = 3
}

variable "instance_storage_space" {
  description = "The storage space of the RabbitMQ instance"
  type        = number
  default     = 600
}

variable "instance_ssl_enable" {
  description = "Whether to enable SSL for the RabbitMQ instance"
  type        = bool
  default     = false
}

variable "instance_access_user_name" {
  description = "The access user of the RabbitMQ instance"
  type        = string
  default     = "admin"
}

variable "instance_password" {
  description = "The access password of the RabbitMQ instance"
  type        = string
  sensitive   = true
  default     = "123456"
  # default     = "YourPassword@123"
}

variable "instance_description" {
  description = "The description of the RabbitMQ instance"
  type        = string
  default     = ""
}

variable "enterprise_project_id" {
  description = "The ID of the enterprise project to which the RabbitMQ instance belongs"
  type        = string
  default     = null
}

variable "instance_tags" {
  description = "The key/value pairs to associate with the RabbitMQ instance"
  type        = map(string)
  default     = {}
}

variable "charging_mode" {
  description = "The charging mode of the RabbitMQ instance"
  type        = string
  default     = "postPaid"
}

variable "period_unit" {
  description = "The period unit of the RabbitMQ instance"
  type        = string
  default     = null
}

variable "period" {
  description = "The period of the RabbitMQ instance"
  type        = number
  default     = null
}

variable "auto_renew" {
  description = "The auto renew of the RabbitMQ instance"
  type        = string
  default     = "false"
}

variable "vhost_name" {
  description = "The name of the RabbitMQ virtual host"
  type        = string
  default     = "app_vhost"
}

variable "exchange_name" {
  description = "The name of the RabbitMQ exchange"
  type        = string
  default     = "app_exchange"
}

variable "exchange_type" {
  description = "The type of the RabbitMQ exchange"
  type        = string
  default     = "direct"
}

variable "queue_name" {
  description = "The name of the RabbitMQ queue"
  type        = string
  default     = "app_queue"
}

variable "producer_instance_name" {
  description = "The name of the producer ECS instance"
  type        = string
}

variable "eip_type" {
  description = "The type of the ECS EIP"
  type        = string
  default     = "5_bgp"
}

variable "eip_share_type" {
  description = "The share type of the ECS EIP"
  type        = string
  default     = "PER"
}

variable "eip_size" {
  description = "The size of the ECS EIP"
  type        = number
  default     = 5
}

variable "eip_charge_mode" {
  description = "The charge mode of the ECS EIP"
  type        = string
  default     = "traffic"
}

variable "message_interval" {
  description = "The interval in seconds between messages sent by the producer"
  type        = number
  default     = 5
}

variable "consumer_instance_name" {
  description = "The name of the consumer ECS instance"
  type        = string
}
