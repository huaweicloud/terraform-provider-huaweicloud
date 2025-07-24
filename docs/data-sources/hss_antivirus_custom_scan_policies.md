---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_custom_scan_policies"
description: |-
  Use this data source to get the list of HSS custom scan policies within HuaweiCloud.
---
# huaweicloud_hss_antivirus_custom_scan_policies

Use this data source to get the list of HSS custom scan policies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_custom_scan_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_name` - (Optional, String) Specifies the policy name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of custom scan policies.

* `data_list` - The list of custom scan policies details.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `policy_id` - The policy ID.

* `policy_name` - The policy name.

* `start_type` - The startup type.  
  The valid values are as follows:
  + **now**
  + **later**
  + **period**

* `scan_period` - The startup period.  
  The valid values are as follows:
  + **day**
  + **week**
  + **month**

* `scan_period_date` - The scan cycle date. The valid values are `1` - `28`.
  + When the `scan_period` is week, `1` - `7` represents sunday to saturday.
  + When the `scan_period` is month, `1` - `28` represents the 1st to 28th of each month.

* `scan_time` - The scan timestamp in milliseconds. Only has a value when the `startup_type` is later.

* `scan_hour` - The scanning hours. Only has a value when the `startup_type` is period.

* `scan_minute` - The scanning minutes. Only has a value when the `startup_type` is period.

* `next_start_time` - The next startup time in milliseconds.

* `scan_dir` - The scan directories. Separate multiple directories with semicolons.

* `ignore_dir` - The exclude directories. Separate multiple directories with semicolons.

* `action` - The disposal action.  
  The valid values are as follows:
  + **auto**: Automatic disposal.
  + **manual**: Manual disposal.

* `invalidate` - Is it invalid. The valid value can be **true** or **false**.

* `host_num` - Affects the number of hosts.

* `host_info_list` - The host details.
  The [host_info_list](#host_info_list_struct) structure is documented below.

* `whether_paid_task` - Is the scanning task paid for this time. The valid value can be **true** or **false**.

* `file_type_list` - The file type list.

<a name="host_info_list_struct"></a>
The `host_info_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `private_ip` - The host private IP address.

* `public_ip` - The host public IP address.

* `asset_value` - The asset importance. The valid value can be **important**, **common**, or **test**.
