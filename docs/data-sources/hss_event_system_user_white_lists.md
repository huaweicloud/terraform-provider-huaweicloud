---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_system_user_white_lists"
description: |-
  Use this data source to get the list of HSS system user white lists within HuaweiCloud.
---
# huaweicloud_hss_event_system_user_white_lists

Use this data source to get the list of HSS system user white lists within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_event_system_user_white_lists" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the private IP of the host.

* `public_ip` - (Optional, String) Specifies the public IP of the host.

* `system_user_name` - (Optional, String) Specifies the system user-name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of white lists.

* `remain_num` - The number of white lists that can continue to be added.

* `limit_num` - The maximum number of white lists.

* `data_list` - The list of white lists details.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `enterprise_project_name` - The enterprise project name.

* `host_id` - The host ID.

* `host_name` - The host name.

* `private_ip` - The private IP address.

* `public_ip` - The public IP address.

* `system_user_name_list` - The list of system user-names.

* `update_time` - The update time in milliseconds.

* `remarks` - The remarks.
