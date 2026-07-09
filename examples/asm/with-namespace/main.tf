# Create a CCE namespace for ASM mesh sidecar injection
resource "huaweicloud_cce_namespace" "test" {
  count = length(var.namespaces) == 0 ? 1 : 0

  cluster_id = var.cluster_id
  name       = var.namespace_name
}

# Create an ASM mesh with namespace configuration
resource "huaweicloud_asm_mesh" "test" {
  name    = var.mesh_name
  type    = var.mesh_type
  version = var.mesh_version
  tags    = var.tags

  extend_params {
    clusters {
      cluster_id = var.cluster_id

      installation {
        nodes {
          field_selector {
            key      = "UID"
            operator = "In"
            values   = [var.node_id]
          }
        }
      }

      injection {
        namespaces {
          field_selector {
            key      = "Name"
            operator = "In"
            values   = length(var.namespaces) > 0 ? var.namespaces : [huaweicloud_cce_namespace.test[0].name]
          }
        }
      }
    }
  }
}
