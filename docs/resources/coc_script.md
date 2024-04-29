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

* `parameters` - (Optional, List) Specifies the input parameters of the script.
  Up to 20 script parameters can be added. The [parameters](#block--parameters) structure is documented below.

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
