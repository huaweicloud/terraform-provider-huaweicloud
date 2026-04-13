---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_organization_tree"
description: |-
  Use this data source to get the organization tree within HuaweiCloud.
---

# huaweicloud_cfw_organization_tree

Use this data source to get the organization tree within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_organization_tree" "test" { 
  fw_instance_id = var.fw_instance_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `parent_id` - (Optional, String) Specifies the parent organization unit ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of organization tree nodes.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `delegated` - The indication of whether the organization unit is delegated.

* `id` - The organization unit ID.

* `name` - The organization unit name.

* `org_type` - The organization unit type.

* `parent_id` - The parent organization unit ID.

* `urn` - The uniform resource name of the organization unit.
