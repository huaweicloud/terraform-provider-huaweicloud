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

# Variables definitions for resource/data source
variable "availability_zone" {
  description = "The availability zone to which the ECS instance belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_id" {
  description = "The flavor ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_flavor_performance_type" {
  description = "The performance type of the ECS instance flavor"
  type        = string
  default     = "normal"
}

variable "instance_flavor_cpu_core_count" {
  description = "The CPU core count of the ECS instance flavor"
  type        = number
  default     = 2
}

variable "instance_flavor_memory_size" {
  description = "The memory size of the ECS instance flavor"
  type        = number
  default     = 4
}

variable "instance_image_id" {
  description = "The image ID of the ECS instance"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_image_visibility" {
  description = "The visibility of the ECS instance image"
  type        = string
  default     = "public"
}

variable "instance_image_os" {
  description = "The OS of the ECS instance image"
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
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
}

variable "subnet_dns_list" {
  description = "The DNS list of the subnet"
  type        = list(string)
  default     = null
}

variable "loadbalancer_name" {
  description = "The name of the load balancer"
  type        = string
}

variable "is_associate_eip" {
  description = "Whether to associate an EIP with the load balancer"
  type        = bool
  default     = false
}

variable "eip_address" {
  description = "The address of the EIP"
  type        = string
  default     = ""
  nullable    = false
}

variable "bandwidth_name" {
  description = "The name of the EIP bandwidth"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = !var.is_associate_eip || (var.eip_address != "" || var.bandwidth_name != "")
    error_message = "The bandwidth name must be provided if the EIP address is not provided."
  }
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

variable "listener_protocol" {
  description = "The protocol of the listener"
  type        = string
  default     = "UDP"
}

variable "listener_default_tls_container_ref" {
  description = "The ID of the server certificate"
  type        = string
  default     = ""
}

variable "listener_server_certificate_name" {
  description = "The name of the server certificate"
  type        = string
  default     = ""

  validation {
    condition     = var.listener_protocol != "TERMINATED_HTTPS" || var.listener_default_tls_container_ref != "" || var.listener_server_certificate_name != ""
    error_message = "listener_server_certificate_name must be provided if listener_protocol is TERMINATED_HTTPS and listener_default_tls_container_ref is not provided"
  }
}

variable "listener_server_certificate_private_key" {
  description = "The private key of the server certificate"
  type        = string
  sensitive   = true
  default     = ""

  validation {
    condition     = var.listener_server_certificate_name == "" || var.listener_server_certificate_private_key != ""
    error_message = "listener_server_certificate_private_key must be provided if listener_server_certificate_name is provided"
  }
}

variable "listener_server_certificate_certificate" {
  description = "The content of the server certificate"
  type        = string
  sensitive   = true
  default     = ""

  validation {
    condition     = var.listener_server_certificate_name == "" || var.listener_server_certificate_certificate != ""
    error_message = "listener_server_certificate_certificate must be provided if listener_server_certificate_name is provided"
  }
}

variable "listener_name" {
  description = "The name of the listener"
  type        = string
}

variable "listener_port" {
  description = "The port of the listener"
  type        = number
  default     = 80
}

variable "listener_description" {
  description = "The description of the listener"
  type        = string
  default     = ""
}

variable "listener_tags" {
  description = "The tags of the listener"
  type        = map(string)
  default     = {}
}

variable "listener_http2_enable" {
  description = "The HTTP/2 enable of the listener, only valid when listener_protocol is TERMINATED_HTTPS"
  type        = bool
  default     = false
}

variable "listener_client_ca_tls_container_ref" {
  description = "The ID of the CA certificate of the listener, only valid when listener_protocol is TERMINATED_HTTPS"
  type        = string
  default     = null
}

variable "listener_sni_container_refs" {
  description = "The ID list of the SNI certificates of the listener, only valid when listener_protocol is TERMINATED_HTTPS"
  type        = list(string)
  default     = null
}

variable "listener_tls_ciphers_policy" {
  description = "The TLS ciphers policy of the listener, only valid when listener_protocol is TERMINATED_HTTPS"
  type        = string
  default     = null
}

variable "listener_insert_headers" {
  description = "The insert headers of the listener, only valid when listener_protocol is TERMINATED_HTTPS"
  type        = object({
    x_forwarded_elb_ip = optional(string, null)
    x_forwarded_host   = optional(string, null)
  })

  default  = {}
  nullable = false
}

variable "pool_name" {
  description = "The name of the backend server group"
  type        = string
  default     = ""
}

variable "pool_protocol" {
  description = "The protocol of the backend server group"
  type        = string
  default     = "UDP"
}

variable "pool_method" {
  description = "The load balancing algorithm of the backend server group"
  type        = string
  default     = "ROUND_ROBIN"
}

variable "pool_description" {
  description = "The description of the backend server group"
  type        = string
  default     = ""
}

variable "pool_persistence" {
  description = "The persistence of the backend server group"
  type        = object({
    type        = string
    cookie_name = optional(string, null)
    timeout     = optional(number, null)
  })

  default = null
}

variable "security_group_name" {
  description = "The name of the security group"
  type        = string
}

variable "instance_name" {
  description = "The name of the ECS instance"
  type        = string
}

variable "member_protocol_port" {
  description = "The protocol port of the backend server"
  type        = number
  default     = 80
}

variable "member_weight" {
  description = "The weight of the backend server"
  type        = number
  default     = 1
}

variable "health_check_port" {
  description = "The port for health checks"
  type        = number
  default     = null
}

variable "security_group_rule_remote_ip_prefix" {
  description = "The remote IP prefix of the security group rule"
  type        = string
  default     = "100.125.0.0/16"
}

variable "health_check_name" {
  description = "The name of the health check"
  type        = string
  default     = "health_check"
}

variable "health_check_type" {
  description = "The type of the health check"
  type        = string
  default     = "UDP_CONNECT"
}

variable "health_check_delay" {
  description = "The delay between health checks in seconds"
  type        = number
  default     = 10
}

variable "health_check_timeout" {
  description = "The timeout for health checks in seconds"
  type        = number
  default     = 5
}

variable "health_check_max_retries" {
  description = "The maximum number of retries for health checks"
  type        = number
  default     = 3
}

variable "health_check_url_path" {
  description = "The URL path for the health check"
  type        = string
  default     = null
}

variable "health_check_http_method" {
  description = "The HTTP method for the health check"
  type        = string
  default     = null
}

variable "health_check_expected_codes" {
  description = "The expected HTTP status codes for the health check"
  type        = string
  default     = null
}

variable "health_check_domain_name" {
  description = "The domain name for the health check"
  type        = string
  default     = null
}
