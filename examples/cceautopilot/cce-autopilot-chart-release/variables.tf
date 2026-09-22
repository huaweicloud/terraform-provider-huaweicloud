# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
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

# Variable definitions for cluster
variable "vpc_name" {
  description = "The VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block of the VPC"
  type        = string
  default     = "192.168.0.0/16"
}

variable "subnet_name" {
  description = "The subnet name"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR block of the subnet"
  type        = string
  default     = ""
}

variable "gateway_ip" {
  description = "The gateway IP address of the subnet"
  type        = string
  default     = ""
}

variable "cluster_name" {
  description = "The CCE autopilot cluster name"
  type        = string
}

variable "cluster_version" {
  description = "The Kubernetes version of the cluster"
  type        = string
  default     = "v1.34"
}

variable "cluster_alias" {
  description = "The alias of the cluster"
  type        = string
  default     = ""
}

variable "cluster_description" {
  description = "The description of the cluster"
  type        = string
  default     = ""
}

variable "cluster_category" {
  description = "The cluster type. Only Turbo is supported"
  type        = string
  default     = "Turbo"
}

variable "cluster_type" {
  description = "The master node architecture. Valid values are: VirtualMachine"
  type        = string
  default     = "VirtualMachine"
}

variable "cluster_custom_san" {
  description = "The custom SAN field in the API server certificate of the cluster"
  type        = list(string)
  default     = []
}

variable "cluster_eip_id" {
  description = "The EIP ID of the cluster"
  type        = string
  default     = ""
}

variable "cluster_enable_snat" {
  description = "Whether SNAT is configured for the cluster"
  type        = bool
  default     = false
}

variable "cluster_enable_swr_image_access" {
  description = "Whether SWR image access is enabled for the cluster"
  type        = bool
  default     = false
}

variable "cluster_delete_efs" {
  description = "Whether to delete the associated EFS when deleting the cluster"
  type        = bool
  default     = false
}

variable "cluster_delete_eni" {
  description = "Whether to delete the associated ENI when deleting the cluster"
  type        = bool
  default     = false
}

variable "cluster_delete_net" {
  description = "Whether to delete the associated network resources when deleting the cluster"
  type        = bool
  default     = false
}

variable "cluster_delete_obs" {
  description = "Whether to delete the associated OBS when deleting the cluster"
  type        = bool
  default     = false
}

variable "cluster_delete_sfs_turbo" {
  description = "Whether to delete the associated SFS Turbo when deleting the cluster"
  type        = bool
  default     = false
}

variable "cluster_lts_reclaim_policy" {
  description = "The LTS reclaim policy. Valid values are: delete, retain"
  type        = string
  default     = "retain"
}

variable "cluster_enterprise_project_id" {
  description = "The ID of the enterprise project to which the cluster belongs"
  type        = string
  default     = "0"
}

variable "cluster_configurations_override" {
  description = "The component configuration items for override"
  type        = list(object({
    name  = string
    value = string
  }))
  default     = []
}

variable "cluster_configurations_override_name" {
  description = "The component name for configurations override"
  type        = string
  default     = ""
}

variable "cluster_tags" {
  description = "The tags of the cluster"
  type        = map(string)
  default     = {}
}

# Variable definitions for chart
variable "chart_content" {
  description = "The path of the chart package to be uploaded"
  type        = string
}

variable "chart_parameters" {
  description = "The parameters of the chart"
  type        = string
  default     = "{\"override\":true,\"skip_lint\":true,\"source\":\"package\"}"
}

# Variable definitions for release
variable "release_name" {
  description = "The name of the release"
  type        = string
}

variable "release_namespace" {
  description = "The namespace to deploy the release"
  type        = string
  default     = "default"
}

variable "release_version" {
  description = "The version of the release"
  type        = string
}

variable "release_description" {
  description = "The description of the release"
  type        = string
  default     = ""
}

variable "release_action" {
  description = "The release updating action. Valid values are: upgrade, rollback"
  type        = string
  default     = ""
}

variable "release_image_tag" {
  description = "The image tag"
  type        = string
  default     = "latest"
}

variable "release_image_pull_policy" {
  description = "The image pull policy"
  type        = string
  default     = "IfNotPresent"
}

variable "release_dry_run" {
  description = "Whether to dry run"
  type        = bool
  default     = false
}

variable "release_name_template" {
  description = "The release name template"
  type        = string
  default     = ""
}

variable "release_no_hooks" {
  description = "Whether to disable hooks during installation"
  type        = bool
  default     = false
}

variable "release_replace" {
  description = "Whether to replace the release with the same name"
  type        = bool
  default     = false
}

variable "release_recreate" {
  description = "Whether to rebuild the release"
  type        = bool
  default     = false
}

variable "release_reset_values" {
  description = "Whether to reset values during an update"
  type        = bool
  default     = false
}

variable "release_rollback_version" {
  description = "The version of the rollback release"
  type        = number
  default     = 0
}

variable "release_include_hooks" {
  description = "Whether to enable hooks during an update or deletion"
  type        = bool
  default     = false
}
