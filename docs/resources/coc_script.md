---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script"
description: ""
---

# huaweicloud_coc_script

Manages a COC script resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_coc_script" "test" {
  name        = "demo"
  description = "a demo script"
  risk_level  = "LOW"
  version     = "1.0.0"
  type        = "SHELL"

  content = <<EOF
#! /bin/bash
echo "hello $${name}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the script. The value can contains 3 to 64 characters,
  including letters, digits, hyphens (-), and underscores (_). Changing this creates a new resource.

* `description` - (Required, String) Specifies the description of the script.
  The value can consist of up to 256 characters.

* `risk_level` - (Required, String) Specifies the risk level. The valid values are **LOW**, **MEDIUM** and **HIGH**.

* `version` - (Required, String) Specifies the version of the script. For example, **1.0.0** or **1.1.0**.

* `type` - (Required, String, ForceNew) Specifies the content type of the script.
  The valid values are **SHELL**, **PYTHON** and **BAT**. Changing this creates a new resource.

* `content` - (Required, String) Specifies the content of the script.
  The value can consist of up to 4096 characters.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID. The default value is **0**.

* `reviewers` - (Optional, List) Specifies the reviewers of the script.
  If left blank, no approval is required. The [reviewers](#block--reviewers) structure is documented below.

* `protocol` - (Optional, String) Specifies the approval message notification protocol, used to notify reviewers.
  Values can be as follows:
  + **DEFAULT**: Default.
  + **SMS**: SMS.
  + **EMAIL**: Email.
  + **DING_TALK**: DingTalk.
  + **WE_LINK**: WeLink.
  + **WECHAT**: WeChat.
  + **CALLNOTIFY**: Language.
  + **NOT_TO_NOTIFY**: Do not notify.

* `parameters` - (Optional, List) Specifies the input parameters of the script.
  Up to 20 script parameters can be added. The [parameters](#block--parameters) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the script.

<a name="block--reviewers"></a>
The `reviewers` block supports:

* `reviewer_name` - (Required, String) Specifies the name of the reviewer.

* `reviewer_id` - (Required, String) Specifies the ID of the reviewer.

<a name="block--parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the name of the parameter.

* `value` - (Required, String) Specifies the **default** value of the parameter.

* `description` - (Required, String) Specifies the description of the parameter.

* `sensitive` - (Optional, Bool) Specifies whether the parameter is sensitive.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
* `status` - The status of the script.
* `created_at` - The creation time of the script.
* `updated_at` - The latest update time of the script.

## Import

The COC script can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_script.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `protocol`.
It is generally recommended running `terraform plan` after importing a script.
You can then decide if changes should be applied to the script, or the resource definition should be updated to align
with the script. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_script" "test" {
    ...

  lifecycle {
    ignore_changes = [
      protocol
    ]
  }
}
```
