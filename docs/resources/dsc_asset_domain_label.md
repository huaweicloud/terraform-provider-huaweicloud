---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_asset_domain_label"
description: |-
  Manages a DSC asset domain label resource within HuaweiCloud.
---

# huaweicloud_dsc_asset_domain_label

Manages a DSC asset domain label resource within HuaweiCloud.

-> If the label has associated assets, deleting the label will move the assets under the label to the ungrouped label.

## Example Usage

### Create a top-level label

```hcl
variable "name" {}

resource "huaweicloud_dsc_asset_domain_label" "test" {
  name      = var.name
  parent_id = "top label"
}
```

### Create a child label

```hcl
variable "name" {}
variable "parent_id" {}

resource "huaweicloud_dsc_asset_domain_label" "test" {
  name      = var.name
  parent_id = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the asset domain label is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the asset domain label.

* `parent_id` - (Required, String, NonUpdatable) Specifies the ID of the parent to which the label belongs.  
  When creating a top-level asset domain label, this parameter value must be set to `top label`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `depth` - The depth of the label.

* `is_leaf` - Whether the label is a leaf node.
  + **1**: The label is a leaf node.
  + **0**: The label is not a leaf node.

## Import

The resource can be imported using the `name` and `parent_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dsc_asset_domain_label.test <name>/<parent_id>
```
