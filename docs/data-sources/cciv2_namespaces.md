---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_namespaces"
description: |-
  Use this data source to get the list of CCI namespaces within HuaweiCloud.
---

# huaweicloud_cciv2_namespaces

Use this data source to get the list of CCI namespaces within HuaweiCloud.

## Example Usage

```hcl
variable "namespace_name" {}

data "huaweicloud_cciv2_namespaces" "test" {
  name = var.namespace_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `namespaces` - The CCI namespaces.
  The [namespaces](#attrblock_namespaces) structure is documented below.

<a name="attrblock_namespaces"></a>
The `namespaces` block supports:

* `annotations` - The annotations of the namespace.

* `api_version` - The API version of the namespace.

* `cluster_name` - The cluster name of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `deletion_grace_period_seconds` - The deletion grace period seconds of the namespace.

* `deletion_timestamp` - The deletion timestamp of the namespace.

* `finalizers` - The finalizers of the namespace.

* `generate_name` - The generate name of the namespace.

* `generation` - The generation of the namespace.

* `kind` - The kind of the namespace.

* `labels` - The labels of the namespace.

* `managed_fields` - The managed fields of the namespace.
  The [managed_fields](#attrblock_namespaces_managed_fields) structure is documented below.

* `owner_references` - The owner references of the namespace.
  The [owner_references](#attrblock_namespaces_owner_references) structure is documented below.

* `resource_version` - The resource version of the namespace.

* `self_link` - The self link of the namespace.

* `status` - The status of the namespace.

* `uid` - The uid of the namespace.

<a name="attrblock_namespaces_managed_fields"></a>
The `managed_fields` block supports:

* `api_version` - The API version of the managed fields.

* `fields_type` - The fields type of the managed fields.

* `fields_v1` - The fields v1 of the managed fields.

* `manager` - The manager of the managed fields.

* `operation` - The operation of the managed fields.

* `time` - The time of the managed fields.

<a name="attrblock_namespaces_owner_references"></a>
The `owner_references` block supports:

* `api_version` - The API version of the owner references.

* `block_owner_deletion` - The block owner deletion of the owner references.

* `controller` - The controller of the owner references.

* `kind` - The kind of the owner references.

* `name` - The name of the owner references.

* `uid` - The uid of the owner references.
