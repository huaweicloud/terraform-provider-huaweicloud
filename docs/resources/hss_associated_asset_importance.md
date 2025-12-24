---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_associated_asset_importance"
description: |-
  Manages an HSS associated asset importance operation resource within HuaweiCloud.
---

# huaweicloud_hss_associated_asset_importance

Manages an HSS associated asset importance operation resource within HuaweiCloud.

-> This resource is a one-time action resource using to operation HSS associated asset importance. Deleting this
  resource will not clear the corresponding request record, but will only remove the resource information from the
  tf state file.

## Example Usage

```hcl
variable "asset_value" {}
variable "host_id_list" {
  type = list(string)
}

resource "huaweicloud_hss_associated_asset_importance" "test" {
  asset_value  = var.asset_value
  host_id_list = var.host_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `asset_value` - (Required, String, NonUpdatable) Specifies the asset importance.  
  The valid values are as follows:
  + **important**: Important assets.
  + **common**: General assets.
  + **test**: Test assets.

* `host_id_list` - (Required, List, NonUpdatable) Specifies the list of host IDs.
  List size ranges from `1` to `200` items.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
