---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_hook"
description: |-
  Manages a RFS private hook resource within HuaweiCloud.
---

# huaweicloud_rfs_private_hook

Manages a RFS private hook resource within HuaweiCloud.

## Example Usage

### Create a private hook using remote OBS object (ZIP file)

```hcl
variable "hook_name" {}
variable "object_access_uri" {}

resource "huaweicloud_rfs_private_hook" "test" {
  name                = var.hook_name
  version             = "1.0.0"
  version_description = "This is a first version"
  policy_uri          = var.object_access_uri

  configuration {
    failure_mode  = "WARN"
    target_stacks = "ALL"
  }
}
```

### Create a private hook using Rego codes

```hcl
variable "hook_name" {}

resource "huaweicloud_rfs_private_hook" "test" {
  name                = var.hook_name
  version             = "1.0.0"
  version_description = "This is a first version"
  policy_body         = <<EOT
package policy

import rego.v1

hook_result := {
  "is_passed": input.message == "world",
  "err_msg": "The error msg when private hook is not passed the validation",
}
EOT

  configuration {
    failure_mode  = "WARN"
    target_stacks = "ALL"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the private hook is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the private hook.  
  The valid length is limited from `1` to `128`, only Chinese or English letters, digits, hyphens (-),
  underscores (_) are allowed.
  The name must start with a Chinese characters or English letter. The names are case sensitive.  
  Change this parameter will create a new resource.

* `version` - (Required, String) Specifies the version of the private hook.
  The version number must follow the **Semantic Version** rules.

* `description` - (Optional, String) Specifies the description of the private hook.  
  The valid length is limited from `1` to `1,024`.

* `version_description` - (Optional, String) Specifies the description of the private hook version.  
  The valid length is limited from `1` to `1,024`.

* `policy_uri` - (Optional, String) Specifies the OBS address of the policy file.  
  The content only supports policy templates written in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language)
  language recognized by the OPA open source engine.
  Policy files currently support single files or zip compressed packages. Single files need to end with `.rego`.
  Compressed packages currently only support zip format and files need to end with `.zip`.
  The verification requirements for policy files are as follows:
  + Size, format, syntax, etc. will be verified when creating.
  + Policy files must be `UTF-8` encoded.
  + The size of a single file or compressed package before and after decompression should be controlled within `1MB`.
  + The number of files in a compressed package cannot exceed `100`.
  + The maximum length of the file path in the compressed package is `2,048`.
  + The maximum length of the file name in the compressed package is `255` bytes.

  -> OBS address supports mutual access between regions of the same type(, regions are divided into general regions and
     dedicated regions. General regions refer to regions that provide general cloud services to public tenants;
     dedicated regions refer to dedicated regions that only carry the same type of business or provide business services
     to specific tenants).

* `policy_body` - (Optional, String) Specifies the policy content of the private hook.  
  Only policy templates written in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language) language that are
  recognized by the OPA open source engine are supported.

-> Exactly one of the `policy_uri` and `policy_body` must be set.

* `configuration` - (Optional, List) Specifies the configuration of the private hook, that can specify the target
  resource stack where the private hook takes effect and the behavior of the resource stack after the private hook
  verification fails.  
  The [configuration](#private_hook_configuration) structure is documented below.

* `keep_old_version` - (Optional, Bool) Specifies whether keeping old version while updating hook version.  
  Defaults to **false**.

  -> A maximum of `199` historical versions can be created for a hook name.

<a name="private_hook_configuration"></a>
The `configuration` block supports:

* `target_stacks` - (Optional, String) Specifies the target resource stack for the private hook to take effect.  
  The valid values are as follows:
  + **NONE**: This private hook will not be applied to any resource stack.
  + **ALL**: This private hook will be applied to all resource stacks under the account.

* `failure_mode` - (Optional, String) Specifies the behavior after private hook verification fails.  
  The valid values are as follows:
  + **FAIL**: The resource stack will stop deploying after this private hook verification fails, and the resource stack
    status will be updated to DEPLOYMENT_FAILED.
  + **WARN**: After this private hook verification fails, only a warning message will be displayed through the resource
    stack event, but it will not affect the resource stack deployment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation of the private hook, in RFC3339 format.

* `updated_at` - The latest update of the private hook, in RFC3339 format.

## Import

Private hooks can be imported using their `name`, e.g.

```bash
$ terraform import huaweicloud_rfs_private_hook.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `policy_uri`, `policy_body` and `keep_old_version`. It is generally
recommended running `terraform plan` after importing a hook. You can keep the resource the same with its definition bo
choosing any of them to update. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rfs_private_hook" "test" {
  ...

  lifecycle {
    ignore_changes = [
      policy_uri,
      keep_old_version,
    ]
  }
}
```
