---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_auto_launch_statistics"
description: |-
  Use this data source to get the list of HSS auto launch statistics within HuaweiCloud.
---

# huaweicloud_hss_auto_launch_statistics

Use this data source to get the list of HSS auto launch statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_auto_launch_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

* `name` - (Optional, String) Specifies the auto launch name.

* `type` - (Optional, String) Specifies the auto launch type. The valid values are as follows:
  + **0**: Auto launch service.
  + **1**: Scheduled task.
  + **2**: Preload dynamic library.
  + **3**: Run registry key.
  + **4**: Startup folder.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The auto launch statistics list.
  The [data_list](#auto_launch_structure) structure is documented below.

<a name="auto_launch_structure"></a>
The `data_list` block supports:

* `name` - The auto launch name.

* `type` - The auto launch type.

* `num` - The number of hosts that have this auto launch item.
