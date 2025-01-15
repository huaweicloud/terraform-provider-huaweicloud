---
subcategory: "Application Service Mesh (ASM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_asm_meshes"
description: |-
  Use this data source to get a list of ASM meshes within HuaweiCloud.
---

# huaweicloud_asm_meshes

Use this data source to get a list of ASM meshes within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_asm_meshes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of meshes.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `name` - The mesh name.
  The name consists of 4 to 64 characters, including letters, digits and hyphens (-),
  must starts with letters and can't end with hyphens (-).

* `id` - The mesh ID.

* `annotations` - The mesh annotations in key/value format.

* `labels` - The mesh labels in key/value format.

* `created_at` - The time when the mesh is created.

* `version` - The mesh version.

* `type` - The mesh type.

* `extend_params` - The extend parameters of the mesh.

  The [extend_params](#spec_extend_params_struct) structure is documented below.

* `tags` - The key/value pairs to associate with the mesh.

* `status` - The status of the mesh.

<a name="spec_extend_params_struct"></a>
The `extend_params` block supports:

* `clusters` - The cluster informations in the mesh.

  The [clusters](#extend_params_clusters_struct) structure is documented below.

<a name="extend_params_clusters_struct"></a>
The `clusters` block supports:

* `cluster_id` - The cluster ID.

* `injection` - The sidecar injection configuration.

  The [injection](#clusters_injection_struct) structure is documented below.

* `installation` - The mesh components installation configuration.

  The [installation](#clusters_installation_struct) structure is documented below.

<a name="clusters_injection_struct"></a>
The `injection` block supports:

* `namespaces` - The namespace of the sidecar injection.

  The [namespaces](#nodes_or_namespaces_struct) structure is documented below.

<a name="clusters_installation_struct"></a>
The `installation` block supports:

* `nodes` - The nodes to install mesh components.

  The [nodes](#nodes_or_namespaces_struct) structure is documented below.

<a name="nodes_or_namespaces_struct"></a>
The `nodes` and `namespaces` block supports:

* `field_selector` - The field selector.

  The [field_selector](#field_selector_struct) structure is documented below.

<a name="field_selector_struct"></a>
The `field_selector` block supports:

* `values` - The value of the selector.

* `key` - The key of the selector.

* `operator` - The operator of the selector.
