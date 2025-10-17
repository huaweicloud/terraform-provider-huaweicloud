---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_whitelist_associate_hosts"
description: |-
  Use this data source to get the list of hosts associated with the process whitelist policy.
---

# huaweicloud_hss_app_whitelist_associate_hosts

Use this data source to the list of hosts associated with the process whitelist policy.

## Example Usage

```hcl
data "huaweicloud_hss_app_whitelist_associate_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_name` - (Optional, String) Specifies the policy name.

* `learning_status` - (Optional, String) Specifies the policy learning status.
  The valid values are as follows:
  + **effecting**: Learning completed, strategy activated.
  + **learned**: Learning completed, pending confirmation.
  + **learning**: Learning in progress.
  + **pause**: Paused.
  + **abnormal**: Learning anomaly.

* `apply_status` - (Optional, String) Specifies the policy application status.
  The valid values are as follows:
  + **true**
  + **false**

* `asset_value` - (Optional, String) Specifies the asset importance.
  The valid values are as follows:
  + **important**
  + **common**
  + **test**

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the private IP address of the host.

* `os_type` - (Optional, String) Specifies the OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `policy_id` - (Optional, String) Specifies the policy ID.

* `public_ip` - (Optional, String) Specifies the public IP address of the host.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of process whitelist policies.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `public_ip` - The public IP address of the host.

* `private_ip` - The private IP address of the host.

* `asset_value` - The asset importance.

* `policy_name` - The policy name.

* `event_num` - The number of events on the host.

* `os_type` - The OS type.

* `learning_status` - The policy learning status.

* `apply_status` - Whether the policy has been applied.

* `intercept` - Whether to enable blocking.

* `policy_id` - The policy ID.

* `policy_type` - The policy type.
