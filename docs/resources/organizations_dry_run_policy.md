---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_dry_run_policy"
description: |-
  Manages an Organizations dry-run policy resource within HuaweiCloud.
---

# huaweicloud_organizations_dry_run_policy

Manages an Organizations dry-run policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "policy_name" {}

resource "huaweicloud_organizations_dry_run_policy" "test" {
  name    = var.policy_name
  type    = "service_control_policy"
  content = jsonencode({
    Version = "5.0",

    Statement = [
      {
        Effect = "Deny",
        Action = []
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the dry-run policy.  
  The name can contain `1` to `64` characters, only Chinese characters, letters, digits, underscore (_), hyphens (-)
  and spaces are allowed and the first and last characters cannot be spaces.

* `content` - (Required, String) Specifies the content of the dry-run policy, in JSON format.

* `type` - (Required, String, NonUpdatable) Specifies the type of the dry-run policy.  
  The valid values are as follows:
  + **service_control_policy**: Service control policy.

* `tags` - (Optional, Map) Specifies the key/value pairs associated with the dry-run policy.

* `description` - (Optional, String) Specifies the description of the dry-run policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - The uniform resource name of the dry-run policy.

* `is_builtin` - Whether the dry-run policy is a built-in policy.

## Import

The dry-run policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_dry_run_policy.test <id>
```
