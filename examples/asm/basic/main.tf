# Create a basic ASM mesh for single cluster
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
    }
  }
}
