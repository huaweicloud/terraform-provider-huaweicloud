# Variable definitions for authentication
variable "region_name" {
  description = "The region where the ASM mesh and CCE cluster are located"
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

# Variable definitions for CCE
variable "namespaces" {
  description = "The list of existing namespace names for sidecar injection"
  type        = list(string)
  default     = []
}

variable "cluster_id" {
  description = "The ID of the CCE cluster to be associated with the ASM mesh"
  type        = string
}

variable "namespace_name" {
  description = "The name of the namespace to create"
  type        = string
  default     = ""

  validation {
    condition     = length(var.namespaces) != 0 || var.namespace_name != ""
    error_message = "The 'namespace_name' is required when 'namespaces' is not specified."
  }
}

# Variable definitions for ASM mesh
variable "mesh_name" {
  description = "The name of the ASM mesh"
  type        = string
}

variable "mesh_type" {
  description = "The type of the ASM mesh"
  type        = string
  default     = "InCluster"
  nullable    = false
}

variable "mesh_version" {
  description = "The version of the ASM mesh"
  type        = string
}

variable "tags" {
  description = "The key/value pairs to associate with the ASM mesh"
  type        = map(string)
  default     = {}
}

variable "node_id" {
  description = "The ID of the node where ASM mesh components will be installed"
  type        = string
}
