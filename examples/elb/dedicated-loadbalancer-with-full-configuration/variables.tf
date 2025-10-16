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

variable "instance_flavor_id" {
  description = "The flavor ID of the instance"
  type        = string
  default     = ""
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

variable "listener_name" {
  description = "The name of the listener"
  type        = string
}

variable "listener_protocol" {
  description = "The protocol of the listener"
  type        = string
  default     = "TCP"
}

variable "listener_port" {
  description = "The port of the listener"
  type        = number
  default     = 8080
}

variable "listener_server_certificate" {
  description = "The server certificate ID of the listener, required when the listener_protocol is HTTPS, TLS or QUIC"
  type        = string
  default     = null
}

variable "listener_ca_certificate" {
  description = "The CA certificate ID of the listener, only available when the listener_protocol is HTTPS"
  type        = string
  default     = null
}

variable "listener_sni_certificates" {
  description = "The SNI certificates of the listener, only available when the listener_protocol is HTTPS or TLS"
  type        = list(string)
  default     = []
}

variable "listener_sni_match_algo" {
  description = "The SNI match algorithm of the listener"
  type        = string
  default     = null
}

variable "listener_security_policy_id" {
  description = "The security policy ID of the listener, only available when the listener_protocol is HTTPS"
  type        = string
  default     = null
}

variable "listener_http2_enable" {
  description = "Whether to enable HTTP/2, only available when the listener_protocol is HTTPS"
  type        = bool
  default     = null
}

variable "listener_port_ranges" {
  description = "The port ranges of the listener"
  type        = list(object({
    start_port = number
    end_port   = number
  }))

  default  = []
  nullable = false
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

variable "listener_advanced_forwarding_enabled" {
  description = "Whether to enable advanced forwarding"
  type        = bool
  default     = false
}

variable "pool_name" {
  description = "The name of the pool"
  type        = string
  default     = null
}

variable "pool_protocol" {
  description = "The protocol of the pool"
  type        = string
  default     = "TCP"
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
  type        = list(object({
    type        = string
    cookie_name = optional(string, null)
    timeout     = optional(number, null)
  }))

  default  = []
  nullable = false
}

variable "health_check_protocol" {
  description = "The protocol of the health check"
  type        = string
  default     = "TCP"
}

variable "health_check_interval" {
  description = "The interval of the health check"
  type        = number
  default     = 20
}

variable "health_check_timeout" {
  description = "The timeout of the health check"
  type        = number
  default     = 15
}

variable "health_check_max_retries" {
  description = "The maximum retries of the health check"
  type        = number
  default     = 10
}

variable "health_check_port" {
  description = "The port of the health check"
  type        = number
  default     = null
}

variable "health_check_url_path" {
  description = "The URL path of the health check"
  type        = string
  default     = null
}

variable "health_check_status_code" {
  description = "The status code of the health check"
  type        = string
  default     = null
}

variable "health_check_http_method" {
  description = "The HTTP method of the health check"
  type        = string
  default     = null
}

variable "health_check_domain_name" {
  description = "The domain name of the health check"
  type        = string
  default     = null
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "member_protocol_port" {
  description = "The port of the member"
  type        = number
  default     = 8080
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "member_weight" {
  description = "The weight of the member"
  type        = number
  default     = 1
}
