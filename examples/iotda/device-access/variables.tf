# Variable definitions for authentication
variable "region_name" {
  description = "The region where the IoTDA resources are located"
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

# Variable definitions for huaweicloud_iotda_space
variable "space_name" {
  description = "The IoTDA resource space name"
  type        = string
}

# Variable definitions for huaweicloud_iotda_product
variable "product_name" {
  description = "The IoTDA product name"
  type        = string
}

variable "product_device_type" {
  description = "The device type of the product"
  type        = string
  default     = "Thermometer"
}

variable "product_protocol" {
  description = "The protocol used by the product"
  type        = string
  default     = "MQTT"
}

variable "product_data_type" {
  description = "The data type of the product"
  type        = string
  default     = "json"
}

variable "product_service_id" {
  description = "The service ID of the product"
  type        = string
}

variable "product_service_type" {
  description = "The service type of the product"
  type        = string
  default     = "serv_type"
}

# Variable definitions for huaweicloud_iotda_device
variable "device_node_id" {
  description = "The node ID of the device"
  type        = string
}

variable "device_name" {
  description = "The device name"
  type        = string
}

variable "device_secret" {
  description = "The secret of the device for authentication"
  type        = string
  sensitive   = true
}

# Variable definitions for IoTDA endpoint
variable "iotda_endpoint" {
  description = "The IoTDA service endpoint for standard/enterprise edition instances"
  type        = string
}
