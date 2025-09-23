---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_storage_classes"
description: |-
  Use this data source to get the list of CCI storage classes within HuaweiCloud.
---

# huaweicloud_cciv2_storage_classes

Use this data source to get the list of CCI storage classes within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cciv2_storage_classes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_classes` - The storage classes.
  The [storage_classes](#storage_classes) structure is documented below.

<a name="storage_classes"></a>
The `storage_classes` block supports:

* `name` - The name.

* `allow_volume_expansion` - The allow volume expansion.

* `allowed_topologies` - The allowed topologies.
  The [allowed_topologies](#storage_classes_allowed_topologies) structure is documented below.

* `annotations` - The annotations.

* `creation_timestamp` - The creation timestamp.

* `labels` - The labels.

* `mount_options` - The mount options.

* `parameters` - The parameters.

* `provisioner` - The provisioner.

* `reclaim_policy` - The reclaim policy.

* `resource_version` - The resource version.

* `uid` - The uid.

* `volume_binding_mode` - The volume binding mode.

<a name="storage_classes_allowed_topologies"></a>
The `allowed_topologies` block supports:

* `match_label_expressions` - The match label expressions.
  The [match_label_expressions](#storage_classes_allowed_topologies_match_label_expressions) structure is documented below.

<a name="storage_classes_allowed_topologies_match_label_expressions"></a>
The `match_label_expressions` block supports:

* `key` - The key.

* `values` - The values.
