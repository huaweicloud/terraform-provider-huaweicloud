---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_batch_config"
description: |-
  Manages a resource to batch configure host management settings within HuaweiCloud.
---

# huaweicloud_hss_host_batch_config

Manages a resource to batch configure host management settings within HuaweiCloud.

-> This resource is a one-time action resource used to batch configure host management Agent resource limits.
  Deleting this resource will not undo the configuration, but will only remove the resource information
  from the tf state file.

## Example Usage

```hcl
variable "operate_all" {}
variable "host_ids" {
  type = list(string)
}
variable "enterprise_project_id" {}
variable "mode" {}
variable "cpu_limit" {}
variable "mem_limit" {}

resource "huaweicloud_hss_host_batch_config" "test" {
  operate_all           = var.operate_all
  enterprise_project_id = var.enterprise_project_id
  host_ids              = var.host_ids
  mode                  = var.mode
  cpu_limit             = var.cpu_limit
  mem_limit             = var.mem_limit
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `operate_all` - (Required, Bool, NonUpdatable) Specifies whether to process all host Agents.
  + **true**: Process all hosts.
  + **false**: Process only the specified hosts (need to set `host_ids`).

* `host_ids` - (Optional, List, NonUpdatable) Specifies the list of host IDs to batch configure.
  Required when `operate_all` is set to **false**. The maximum number of host IDs is `200`.

* `mode` - (Required, String, NonUpdatable) Specifies the resource limit type.
  The valid values are as follows:
  + **default**: Default rule.
  + **customized**: User-defined rule.
  + **adaptive**: Adaptive rule.

* `cpu_limit` - (Required, String, NonUpdatable) Specifies the maximum CPU limit.
  The length is `0` to `32` characters.

* `mem_limit` - (Required, String, NonUpdatable) Specifies the maximum memory limit.
  The length is `0` to `32` characters.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is required only when the enterprise project feature is enabled.
  If you want to query assets under all enterprise projects, set this parameter to **all_granted_eps**.
  Defaults to **0** (the default enterprise project).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
