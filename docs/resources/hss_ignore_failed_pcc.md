---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ignore_failed_pcc"
description: |-
  Using this resource to ignore or cancel the ignorance of servers that fail the password complexity check within HuaweiCloud.
---

# huaweicloud_hss_ignore_failed_pcc

Using this resource to ignore or cancel the ignorance of servers that fail the password complexity check within HuaweiCloud.

-> This resource is a stateless operation resource. Deleting this resource will not affect the ignore status of the hosts,
but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_hss_ignore_failed_pcc" "test" {
  action      = "ignore"
  operate_all = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the action type to perform.
  The valid values are:
  + **ignore**: Ignore the servers that fail the password complexity check.
  + **unignore**: Unignore the servers that fail the password complexity check.

* `operate_all` - (Optional, Bool, NonUpdatable) Specifies whether the operation is a full operation. A maximum of
  `1,000` hosts can be processed at a time. Defaults to **false**.

* `host_ids` - (Optional, List, NonUpdatable) Specifies the list of host IDs to perform the action on.
  This parameter is ignored when `operate_all` is set to **true**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project that the server
  belongs to. The value **0** indicates the default enterprise project. To query servers in all enterprise projects,
  set this parameter to **all_granted_eps**. If you have only the permission on an enterprise project, you need to
  transfer the enterprise project ID to query the server in the enterprise project. Otherwise, an error is reported due
  to insufficient permission.

  -> An enterprise project can be configured only after the enterprise project function is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
