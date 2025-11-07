---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_protect_switch_mode"
description: |-
  Manages an HSS switch cluster protection mode resource within HuaweiCloud.
---

# huaweicloud_hss_cluster_protect_switch_mode

Manages an HSS switch cluster protection mode resource within HuaweiCloud.

-> This resource is only a one-time action resource used to switch HSS cluster protect mode. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "cluster_ids" {
  type = list(string)
}
variable "opr" {}

resource "huaweicloud_hss_cluster_protect_switch_mode" "test" {
  cluster_ids = var.cluster_ids
  opr         = var.opr
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_ids` - (Required, List, NonUpdatable) Specifies the cluster ID list.  

* `opr` - (Required, Int, NonUpdatable) Specifies the switch mode.  
  The valid values are as follows:
  + **1**: Open.
  + **0**: Close.
  + **2**: Close and uninstall plugins.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
