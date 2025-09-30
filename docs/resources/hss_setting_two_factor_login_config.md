---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_two_factor_login_config"
description: |-
  Manages a resource to set two-factor login configuration within HuaweiCloud.
---

# huaweicloud_hss_setting_two_factor_login_config

Manages a resource to set two-factor login configuration within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "host_id_list" {
  type = list(string)
}

resource "huaweicloud_smn_topic" "test" {
  ...
}

resource "huaweicloud_hss_setting_two_factor_login_config" "test" {
  enabled            = true
  auth_type          = "sms"
  host_id_list       = var.host_id_list
  topic_display_name = huaweicloud_smn_topic.test.display_name
  topic_urn          = huaweicloud_smn_topic.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `enabled` - (Required, Bool, NonUpdatable) Specifies whether the two-factor authentication enabled.

* `auth_type` - (Required, String, NonUpdatable) Specifies the authentication type.
  The valid values are as follows:
  + **sms**: Indicates SMS and email verification.
  + **code**: Indicates captcha code verification.

* `host_id_list` - (Required, List, NonUpdatable) Specifies the host IDs.

* `topic_display_name` - (Optional, String, NonUpdatable) Specifies the SMN topic display name.

* `topic_urn` - (Optional, String, NonUpdatable) Specifies the SMN topic urn.

-> The field `topic_display_name` and `topic_urn` are valid only when `auth_type` is set to **sms**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
