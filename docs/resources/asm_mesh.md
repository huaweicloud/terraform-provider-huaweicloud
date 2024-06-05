---
subcategory: "Application Service Mesh (ASM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_asm_mesh"
description: |-
  Manages a ASM mesh resource within HuaweiCloud.
---

# huaweicloud_asm_mesh

Manages a ASM mesh resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "cluster_id" {}
variable "node_id" {}

resource "huaweicloud_asm_mesh" "test" {
  name    = var.name
  type    = "InCluster"
  version = "1.18.7-r1"

  tags = {
    foo = "bar"
    key = "value"
  }

  extend_params {
    clusters {
      cluster_id = var.cluster_id
      installation {
        nodes {
          field_selector {
            key      = "UID"
            operator = "In"
            values   = [
              var.node_id
            ]
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `name` - (Required, String, NonUpdatable) Specifies mesh name.
  The name consists of 4 to 64 characters, including letters, digits and hyphens (-),
  must starts with letters and can't end with hyphens (-).

* `type` - (Required, String, NonUpdatable) Specifies the mesh type.
  The value can be **InCluster**.

* `version` - (Required, String, NonUpdatable) Specifies the mesh version.

* `extend_params` - (Required, List, NonUpdatable) Specifies the extend parameters of the mesh.

  The [extend_params](#spec_extend_params_struct) structure is documented below.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the mesh annotations in key/value format.

* `labels` - (Optional, Map, NonUpdatable) Specifies the mesh labels in key/value format.

* `tags` - (Optional, Map, NonUpdatable) Specifies the key/value pairs to associate with the mesh.

<a name="spec_extend_params_struct"></a>
The `extend_params` block supports:

* `clusters` - (Required, List, NonUpdatable) Specifies the cluster informations in the mesh.

  The [clusters](#extend_params_clusters_struct) structure is documented below.

<a name="extend_params_clusters_struct"></a>
The `clusters` block supports:

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `installation` - (Required, List, NonUpdatable) Specifies the mesh components installation configuration.

  The [installation](#clusters_installation_struct) structure is documented below.

* `injection` - (Optional, List, NonUpdatable) Specifies the sidecar injection configuration.

  The [injection](#clusters_injection_struct) structure is documented below.

<a name="clusters_installation_struct"></a>
The `installation` block supports:

* `nodes` - (Required, List, NonUpdatable) Specifies the mesh components installation configuration.

  The [nodes](#nodes_or_namespaces_struct) structure is documented below.

<a name="clusters_injection_struct"></a>
The `injection` block supports:

* `namespaces` - (Required, List, NonUpdatable) Specifies the namespace of the sidecar injection.

  The [namespaces](#nodes_or_namespaces_struct) structure is documented below.

<a name="nodes_or_namespaces_struct"></a>
The `namespaces` and `nodes` block support:

* `field_selector` - (Required, List, NonUpdatable) Specifies the field selector.

  The [field_selector](#field_selector_struct) structure is documented below.

<a name="field_selector_struct"></a>
The `field_selector` block supports:

* `key` - (Required, String, NonUpdatable) Specifies the key of the selector.

* `operator` - (Required, String, NonUpdatable) Specifies the operator of the selector.
  The value can be **In**.

* `values` - (Required, List, NonUpdatable) Specifies the value of the selector.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The time when the mesh is created.

* `status` - The status of the mesh.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The ASM mesh can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_asm_mesh.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `extend_params`, `annotations`, `labels` and `tags`.
It is generally recommended running `terraform plan` after importing a mesh.
You can then decide if changes should be applied to the mesh, or the resource definition should be updated to
align with the mesh. Also you can ignore changes as below.

```hcl
resource "huaweicloud_asm_mesh" "test" {
    ...

  lifecycle {
    ignore_changes = [
      extend_params, annotations, labels, tags,
    ]
  }
}
```
