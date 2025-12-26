# Variable definitions for authentication
variable "region_name" {
  description = "The region where the CBH instance is located"
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
  description = "The availability zone where the resources will be created"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_id" {
  description = "The ID of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.vpc_id != "" || var.vpc_name != ""
    error_message = "vpc_name must be provided if vpc_id is not provided."
  }
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
  default     = ""
  nullable    = false

  validation {
    condition     = var.subnet_id == "" || var.subnet_name == ""
    error_message = "subnet_name must be provided if subnet_id is not provided."
  }
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "subnet_gateway_ip" {
  description = "The gateway IP of the subnet"
  type        = string
  default     = ""
  nullable    = false
}

variable "eip_address" {
  description = "The EIP address of the CCE cluster"
  type        = string
  default     = ""
  nullable    = false
}

variable "eip_type" {
  description = "The type of the EIP"
  type        = string
  default     = "5_bgp"
}

variable "bandwidth_name" {
  description = "The name of the bandwidth"
  type        = string
  default     = ""
}

variable "bandwidth_size" {
  description = "The size of the bandwidth"
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

variable "cluster_name" {
  description = "The name of the CCE cluster"
  type        = string
}

variable "cluster_flavor_id" {
  description = "The flavor ID of the CCE cluster"
  type        = string
  default     = "cce.s1.small"
}

variable "cluster_version" {
  description = "The version of the CCE cluster"
  type        = string
  default     = null # Default to the latest version
  nullable    = true
}

variable "cluster_type" {
  description = "The type of the CCE cluster"
  type        = string
  default     = "VirtualMachine"
}

variable "container_network_type" {
  description = "The type of container network"
  type        = string
  default     = "overlay_l2"
}

variable "authentication_mode" {
  description = "The mode of the CCE cluster"
  type        = string
  default     = "rbac"
}

variable "delete_all_resources_on_terminal" {
  description = "Whether to delete all resources on terminal"
  type        = bool
  default     = true
}

variable "enterprise_project_id" {
  description = "The enterprise project ID of the CCE cluster"
  type        = string
  default     = "0"
}

variable "node_flavor_id" {
  description = "The flavor ID of the node"
  type        = string
  default     = ""
  nullable    = false
}

variable "node_performance_type" {
  description = "The performance type of the node"
  type        = string
  default     = "general"
}

variable "node_cpu_core_count" {
  description = "The CPU core count of the node"
  type        = number
  default     = 4
}

variable "node_memory_size" {
  description = "The memory size of the node"
  type        = number
  default     = 8
}

variable "keypair_name" {
  description = "The name of the keypair"
  type        = string
}

variable "node_name" {
  description = "The name of the node"
  type        = string
}

variable "root_volume_type" {
  description = "The type of the root volume"
  type        = string
  default     = "SATA"
}

variable "root_volume_size" {
  description = "The size of the root volume"
  type        = number
  default     = 40
}

variable "data_volumes_configuration" {
  description = "The configuration of the data volumes"

  type = list(object({
    volumetype = string
    size       = number
  }))

  default  = []
  nullable = false
}

variable "bucket_name" {
  description = "The name of the OBS bucket"
  type        = string
}

variable "bucket_multi_az" {
  description = "Whether to enable multi-AZ for the OBS bucket"
  type        = bool
  default     = true
}

variable "secret_name" {
  description = "The name of the Kubernetes secret"
  type        = string
}

variable "namespace_name" {
  description = "The name of the CCE namespace"
  type        = string
  default     = "default"
}

variable "secret_labels" {
  description = "The labels of the Kubernetes secret"
  type        = map(string)
  default     = {
    "secret.kubernetes.io/used-by" = "csi"
  }
}

variable "secret_data" {
  description = "The data of the Kubernetes secret"
  type        = map(string)
  sensitive   = true
}

variable "secret_type" {
  description = "The type of the Kubernetes secret"
  type        = string
  default     = "cfe/secure-opaque"
}

variable "pv_name" {
  description = "The name of the persistent volume"
  type        = string
}

variable "pv_csi_provisioner_identity" {
  description = "The identity of the CSI provisioner"
  type        = string
  default     = "everest-csi-provisioner"
}

variable "pv_access_modes" {
  description = "The access modes of the persistent volume"
  type        = list(string)
  default     = ["ReadWriteMany"]
}

variable "pv_storage" {
  description = "The storage of the persistent volume"
  type        = string
  default     = "1Gi"
}

variable "pv_driver" {
  description = "The driver of the persistent volume"
  type        = string
  default     = "obs.csi.everest.io"
}

variable "pv_fstype" {
  description = "The instance type of the persistent volume"
  type        = string
  default     = "s3fs"
}

variable "pv_obs_volume_type" {
  description = "The type of the OBS volume of the persistent volume"
  type        = string
  default     = "standard"
}

variable "pv_reclaim_policy" {
  description = "The reclaim policy of the persistent volume"
  type        = string
  default     = "Retain"
}

variable "pv_storage_class_name" {
  description = "The name of the storage class of the persistent volume"
  type        = string
  default     = "csi-obs"
}

variable "pvc_name" {
  description = "The name of the persistent volume claim"
  type        = string
}

variable "deployment_name" {
  description = "The name of the deployment"
  type        = string
}

variable "deployment_replicas" {
  description = "The number of pods for the deployment"
  type        = number
  default     = 2
}

variable "deployment_containers" {
  description = "The container list for the deployment"
  type        = list(object({
    name          = string
    image         = string
    volume_mounts = list(object({
      mount_path = string
    }))
  }))
  nullable    = false
}

variable "deployment_volume_name" {
  description = "The name of the volume of the deployment"
  type        = string
  default     = "pvc-obs-volume"
}

variable "deployment_image_pull_secrets" {
  description = "The image pull secrets of the deployment"
  type        = list(object({
    name = string
  }))

  default  = [
    {
      name = "default-secret"
    }
  ]
  nullable = false
}
