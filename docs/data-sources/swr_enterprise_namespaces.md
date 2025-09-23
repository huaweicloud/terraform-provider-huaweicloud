---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_namespaces"
description: |-
  Use this data source to get the list of SWR enterprise instance namespaces.
---

# huaweicloud_swr_enterprise_namespaces

Use this data source to get the list of SWR enterprise instance namespaces.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_namespaces" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `name` - (Optional, String) Specifies the namespace name.

* `order_column` - (Optional, String) Specifies the order column.

* `order_type` - (Optional, String) Specifies the order type.
  Value can be **desc** or **asc**, should use with `order_column`.

* `public` - (Optional, String) Specifies whether the namespace is public. Value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `namespaces` - Indicates the namespaces.
  The [namespaces](#attrblock--namespaces) structure is documented below.

<a name="attrblock--namespaces"></a>
The `namespaces` block supports:

* `id` - Indicates the namespace ID.

* `name` - Indicates the namespace name.

* `repo_count` - Indicates the repo count of the namespace.

* `metadata` - Indicates the metadata.
  The [metadata](#attrblock--namespaces--metadata) structure is documented below.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

<a name="attrblock--namespaces--metadata"></a>
The `metadata` block supports:

* `public` - Indicates whether the namespace is public.
