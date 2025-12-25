# Variables definitions for authorization
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

# Variable definitions for resources/data sources
variable "availability_zone" {
  description = "The name of the availability zone to which the resources belong"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "172.16.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "loadbalancer_name" {
  description = "The name of the loadbalancer"
  type        = string
}

variable "loadbalancer_cross_vpc_backend" {
  description = "Whether to associate backend servers with the load balancer by using their IP addresses"
  type        = bool
  default     = false
}

variable "loadbalancer_description" {
  description = "The description of the loadbalancer"
  type        = string
  default     = null
}

variable "enterprise_project_id" {
  description = "The enterprise project ID"
  type        = string
  default     = null
}

variable "loadbalancer_tags" {
  description = "The tags of the loadbalancer"
  type        = map(string)
  default     = {}
}

variable "loadbalancer_force_delete" {
  description = "Whether to force delete the loadbalancer"
  type        = bool
  default     = true
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
  default     = "tf_test_eip"
}

variable "bandwidth_size" {
  description = "The bandwidth size of the EIP in Mbps"
  type        = number
  default     = 5
}

variable "bandwidth_share_type" {
  description = "The share type of the bandwidth"
  type        = string
  default     = "PER"
}

variable "bandwidth_charge_mode" {
  description = "The charge mode of the bandwidth"
  type        = string
  default     = "traffic"
}

variable "listener_name" {
  description = "The name of the listener"
  type        = string
}

variable "listener_protocol" {
  description = "The protocol of the listener"
  type        = string
  default     = "HTTP"
}

variable "listener_port" {
  description = "The port of the listener"
  type        = number
  default     = 8080
}

variable "listener_idle_timeout" {
  description = "The idle timeout of the listener"
  type        = number
  default     = 60
}

variable "listener_request_timeout" {
  description = "The request timeout of the listener"
  type        = number
  default     = null
}

variable "listener_response_timeout" {
  description = "The response timeout of the listener"
  type        = number
  default     = null
}

variable "listener_description" {
  description = "The description of the listener"
  type        = string
  default     = null
}

variable "listener_tags" {
  description = "The tags of the listener"
  type        = map(string)
  default     = {}
}

variable "pool_name" {
  description = "The name of the pool"
  type        = string
  default     = null
}

variable "pool_protocol" {
  description = "The protocol of the pool"
  type        = string
  default     = "HTTP"
}

variable "pool_method" {
  description = "The load balancing method of the pool"
  type        = string
  default     = "ROUND_ROBIN"
}

variable "pool_any_port_enable" {
  description = "Whether to enable any port for the pool"
  type        = bool
  default     = false
}

variable "pool_description" {
  description = "The description of the pool"
  type        = string
  default     = null
}

variable "pool_persistences" {
  description = "The persistence configurations for the pool"

  type = list(object({
    type        = string
    cookie_name = optional(string, null)
    timeout     = optional(number, null)
  }))

  default  = []
  nullable = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the instance flavor"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_image_visibility" {
  description = "The visibility of the instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The OS of the instance image"
  type        = string
  default     = "Ubuntu"
}

variable "configuration_name" {
  description = "The name of the scaling configuration"
  type        = string
}

variable "configuration_image_id" {
  description = "The image ID of the scaling configuration"
  type        = string
  default     = ""
  nullable    = false
}

variable "configuration_flavor_id" {
  description = "The flavor ID of the scaling configuration"
  type        = string
  default     = ""
  nullable    = false
}

variable "configuration_user_data" {
  description = "The user data for the scaling configuration instances"
  type        = string
}

variable "configuration_disks" {
  description = "The disk configurations for the scaling configuration instances"

  type = list(object({
    size        = number
    volume_type = string
    disk_type   = string
  }))

  nullable = false
}

variable "group_name" {
  description = "The name of the AS group"
  type        = string
}

variable "group_desire_instance_number" {
  description = "The desire instance number of the AS group"
  type        = number
  default     = 0
}

variable "group_min_instance_number" {
  description = "The min instance number of the AS group"
  type        = number
  default     = 0
}

variable "group_max_instance_number" {
  description = "The max instance number of the AS group"
  type        = number
  default     = 10
}

variable "group_delete_publicip" {
  description = "Whether to delete the public IP address of the AS group"
  type        = bool
  default     = true
}

variable "group_delete_instances" {
  description = "Whether to delete the instances of the AS group"
  type        = bool
  default     = true
}

variable "group_force_delete" {
  description = "Whether to force delete the AS group"
  type        = bool
  default     = true
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "alarm_rule_name" {
  description = "The name of the CES alarm rule"
  type        = string
}

variable "alarm_rule_conditions" {
  description = "The conditions of the CES alarm rule"

  type = list(object({
    period              = number
    filter              = string
    comparison_operator = string
    value               = number
    unit                = string
    count               = number
    alarm_level         = number
    metric_name         = string
  }))

  nullable = false
}

variable "policy_name" {
  description = "The name of the AS policy"
  type        = string
}

variable "policy_cool_down_time" {
  description = "The cool down time of the AS policy"
  type        = number
  default     = 900
}

variable "policy_operation" {
  description = "The operation of the AS policy"
  type        = string
  default     = "ADD"
}

variable "policy_instance_number" {
  description = "The instance number of the AS policy"
  type        = number
  default     = 1
}
