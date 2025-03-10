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

* `creation_timestamp` - The creation timestamp of the namespace.

* `finalizers` - The finalizers of the namespace.

* `kind` - The kind of the namespace.

* `labels` - The labels of the namespace.

* `resource_version` - The resource version of the namespace.

* `status` - The status of the namespace.

* `uid` - The uid of the namespace.
