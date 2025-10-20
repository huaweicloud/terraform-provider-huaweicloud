---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_whitelist_policy_process_extend"
description: |-
  Use this data source to get the process extend list of HSS app whitelist policy within HuaweiCloud.
---

# huaweicloud_hss_app_whitelist_policy_process_extend

Use this data source to get the process extend list of HSS app whitelist policy within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}
variable "host_id" {}

data "huaweicloud_hss_app_whitelist_policy_process_extend" "test" {
  policy_id = var.policy_id
  host_id   = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the policy ID.

* `host_id` - (Required, String) Specifies the host ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of process extend information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `process_name` - The process name.

* `process_path` - The process path.

* `process_hash` - The process hash.

* `container_id` - The container ID.

* `cmdline` - The process command line.

* `file_size` - The file size.
