---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_two_factor_login_hosts"
description: |-
  Use this data source to get the list of two-factor hosts.
---

# huaweicloud_hss_setting_two_factor_login_hosts

Use this data source to get the list of two-factor hosts.

## Example Usage

```hcl
data "huaweicloud_hss_setting_two_factor_login_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_name` - (Optional, String) Specifies the host name.

* `display_name` - (Optional, String) Specifies the SMN topic display name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The two-factor host list.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `os_type` - The operating system type.
  The valid values are as follows:
  + **linux**
  + **windows**

* `auth_switch` - Whether the two-factor authentication is enabled.

* `auth_type` - The two-factor authentication type.
  The valid values are as follows:
  + **sms**: Indicates SMS and email verification.
  + **code**: Indicates captcha code verification.

* `topic_display_name` - The SMN topic display name.

* `topic_urn` - The SMN topic urn.

* `outside_host` - Whether the host is an external (non-HuaweiCloud) matchine.
