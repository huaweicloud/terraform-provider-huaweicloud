---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_unblock_ip"
description: |-
  Use this data source to get the list of HSS unblock IP within HuaweiCloud.
---

# huaweicloud_hss_event_unblock_ip

Use this data source to get the list of HSS unblock IP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_event_unblock_ip" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `last_days` - (Optional, Int) Specifies the query time range and number of days.

* `host_name` - (Optional, String) Specifies the host name.

* `src_ip` - (Optional, String) Specifies the IP address of the attack source.

* `intercept_status` - (Optional, String) Specifies interception status
  The valid values are as follows:
  + **intercepted**: Indicates that it has been intercepted.
  + **canceled**: Indicates that it has been unblocked.
  + **cancelling**: Indicates pending unblock.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts
  belong.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The details of intercepted IP list.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `src_ip` - The IP address of the attack source.

* `login_type` - The login type.
  The valid values are as follows:
  + **mysql**: Represents the MySQL service.
  + **rdp**: Represents the RDP service.
  + **ssh**: Represents the SSH service.
  + **vsftp**: Represents the VSFTP service.

* `intercept_num` - The number of interceptions.

* `intercept_status` - The interception status.

* `block_time` - The start interception time in milliseconds.

* `latest_time` - The most recent interception time in milliseconds.
