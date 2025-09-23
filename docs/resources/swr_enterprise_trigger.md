---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_trigger"
description: |-
  Manages a SWR enterprise instance trigger resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_trigger

Manages a SWR enterprise instance trigger resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "name" {}

resource "huaweicloud_swr_enterprise_trigger" "test" {
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  name           = var.name
  description    = "desc"
  enabled        = true
  event_types    = ["PUSH_ARTIFACT"]
  
  scope_rules {
    repo_scope_mode = "regular"

    scope_selectors {
      key = "repository"

      value {
        decoration = "repoMatches"
        kind       = "doublestar"
        pattern    = "**"
      }
    }

    tag_selectors {
      decoration = "matches"
      kind       = "doublestar"
      pattern    = "**"
    }
  }

  targets {
    address          = "https://test.com"
    address_type     = "internal"
    auth_header      = "Test:Header"
    skip_cert_verify = false
    type             = "http"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String, NonUpdatable) Specifies the namespace name.

* `name` - (Required, String) Specifies the trigger name.

* `event_types` - (Required, List) Specifies the event types of trigger. Value can be **PUSH_ARTIFACT**.

* `scope_rules` - (Required, List) Specifies the scope rules
  The [scope_rules](#block--scope_rules) structure is documented below.

* `targets` - (Required, List) Specifies the target params.
  The [targets](#block--targets) structure is documented below.

* `description` - (Optional, String) Specifies the description of trigger.

* `enabled` - (Optional, Bool) Specifies whether the trigger is enabled. Default to **false**.

<a name="block--scope_rules"></a>
The `scope_rules` block supports:

* `repo_scope_mode` - (Required, String) Specifies the repository select method. Value can be **regular**.

* `tag_selectors` - (Required, List) Specifies the repository version selector.
  The [tag_selectors](#block--scope_rules--tag_selectors) structure is documented below.

* `scope_selectors` - (Optional, List) Specifies the repository selectors.
  The [scope_selectors](#block--scope_rules--scope_selectors) structure is documented below.

<a name="block--scope_rules--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - (Required, String) Specifies the selector matching type.

* `kind` - (Required, String) Specifies the matching type.

* `pattern` - (Required, String) Specifies the pattern.

<a name="block--scope_rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - (Required, String) Specifies the repository selector key.

* `value` - (Required, List) Specifies the repository selector value.
  The [value](#block--scope_rules--scope_selectors--value) structure is documented below.

<a name="block--scope_rules--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - (Required, String) Specifies the selector matching type.

* `kind` - (Required, String) Specifies the matching type.

* `pattern` - (Required, String) Specifies the pattern.

<a name="block--targets"></a>
The `targets` block supports:

* `address` - (Required, String, NonUpdatable) Specifies the trigger address.

* `address_type` - (Required, String, NonUpdatable) Specifies the trigger address type.

* `type` - (Required, String, NonUpdatable) Specifies the trigger type.

* `auth_header` - (Optional, String) Specifies the auth header. Format is **Key:Value;Key1:Value2**.

* `skip_cert_verify` - (Optional, Bool) Specifies whether to skip the verification of the certificate. Default to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time.

* `creator` - Indicates the creator

* `namespace_id` - Indicates the namespace ID

* `trigger_id` - Indicates the trigger ID.

* `updated_at` - Indicates the last update time.

## Import

The trigger can be imported using `instance_id`, `namespace_name` and `trigger_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_trigger.test <instance_id>/<namespace_name>/<trigger_id>
```
