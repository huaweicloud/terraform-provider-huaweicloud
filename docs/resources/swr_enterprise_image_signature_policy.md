---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_image_signature_policy"
description: |-
  Manages a SWR enterprise image signature policy resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_image_signature_policy

Manages a SWR enterprise image signature policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "name" {}
variable "signature_algorithm" {}
variable "signature_key" {}

resource "huaweicloud_swr_enterprise_image_signature_policy" "test" {
  instance_id         = var.instance_id
  namespace_name      = var.namespace_name
  name                = var.name
  enabled             = true
  signature_method    = "KMS"
  signature_algorithm = var.signature_algorithm
  signature_key       = var.signature_key
  
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
      extras     = jsonencode({
        "untagged": true
      })
    }
  }

  trigger {
    type = "event_based"
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

* `name` - (Required, String) Specifies the policy name.

* `signature_method` - (Required, String) Specifies the signature method. Value can be **KMS**.

* `signature_algorithm` - (Required, String) Specifies the signature algorithm.
  + When the KMS key algorithm is **EC_P256**, value can be **ECDSA_SHA_256**.
  + When the KMS key algorithm is **EC_P384**, value can be **ECDSA_SHA_384**.
  + When the KMS key algorithm is **SM2**, value can be **SM2DSA_SM3**.

* `signature_key` - (Required, String) Specifies the signature key.
  Value can be as follows:
  + **manual**: manual
  + **event_based**: manual and event_based

* `scope_rules` - (Required, List) Specifies the scope rules
  The [scope_rules](#block--scope_rules) structure is documented below.

* `trigger` - (Required, List) Specifies the trigger config.
  The [trigger](#block--trigger) structure is documented below.

* `description` - (Optional, String) Specifies the description of policy.

* `enabled` - (Optional, Bool) Specifies whether the policy is enabled. Default to **false**.

<a name="block--scope_rules"></a>
The `scope_rules` block supports:

* `repo_scope_mode` - (Required, String) Specifies the repository select method.

* `tag_selectors` - (Required, List) Specifies the repository version selector.
  The [tag_selectors](#block--scope_rules--tag_selectors) structure is documented below.

* `scope_selectors` - (Required, List) Specifies the repository selectors.
  The [scope_selectors](#block--scope_rules--scope_selectors) structure is documented below.

<a name="block--scope_rules--tag_selectors"></a>
The `tag_selectors` block supports:

* `decoration` - (Required, String) Specifies the selector matching type. Value can be **repoMatches**.

* `kind` - (Required, String) Specifies the matching type. Value can be **doublestar**.

* `pattern` - (Required, String) Specifies the pattern.

* `extras` - (Optional, String) Specifies the extra infos.

<a name="block--scope_rules--scope_selectors"></a>
The `scope_selectors` block supports:

* `key` - (Required, String) Specifies the repository selector key.

* `value` - (Required, List) Specifies the repository selector value.
  The [value](#block--scope_rules--scope_selectors--value) structure is documented below.

<a name="block--scope_rules--scope_selectors--value"></a>
The `value` block supports:

* `decoration` - (Required, String) Specifies the selector matching type. Value can be **repoMatches**.

* `kind` - (Required, String) Specifies the matching type. Value can be **doublestar**.

* `pattern` - (Required, String) Specifies the pattern.

* `extras` - (Optional, String) Specifies the extra infos.

<a name="block--trigger"></a>
The `trigger` block supports:

* `type` - (Required, String) Specifies the trigger type.

* `trigger_settings` - (Optional, List) Specifies the trigger settings.
  The [trigger_settings](#block--trigger--trigger_settings) structure is documented below.

<a name="block--trigger--trigger_settings"></a>
The `trigger_settings` block supports:

* `cron` - (Optional, String) Specifies the scheduled setting.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `namespace_id` - Indicates the namespace ID

* `policy_id` - Indicates the policy ID.

* `creator` - Indicates the creator

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Import

The image signature policy can be imported using `instance_id`, `namespace_name` and `policy_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_image_signature_policy.test <instance_id>/<namespace_name>/<policy_id>
```
