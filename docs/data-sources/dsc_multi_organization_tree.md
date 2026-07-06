---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_multi_organization_tree"
description: |-
  Use this data source to get the organizational unit tree of an account within HuaweiCloud.
---

# huaweicloud_dsc_multi_organization_tree

Use this data source to get the organizational unit tree of an account within HuaweiCloud.

## Example Usage

```hcl
variable "entity_id" {}

data "huaweicloud_dsc_multi_organization_tree" "test" {
  entity_id = var.entity_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `entity_id` - (Required, String) Specifies the entity ID used to identify a specific account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ou_list` - The organizational unit list.

  The [ou_list](#ou_list_struct) structure is documented below.

<a name="ou_list_struct"></a>
The `ou_list` block supports:

* `id` - The organizational unit ID.

* `name` - The organizational unit name.
