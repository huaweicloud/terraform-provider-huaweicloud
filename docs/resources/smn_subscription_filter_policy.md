---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_subscription_filter_policy"
description: |-
  Manages an SMN subscription filter policy resource within HuaweiCloud.
---

# huaweicloud_smn_subscription_filter_policy

Manages an SMN subscription filter policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "subscription_urn" {}

resource "huaweicloud_smn_subscription_filter_policy" "test" {
  subscription_urn = var.subscription_urn

  filter_policies {
    name          = "alarm"
    string_equals = ["os", "process"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `subscription_urn` - (Required, String, NonUpdatable) Specifies the resource identifier of the subscriber.

* `filter_policies` - (Required, List) Specifies the message filter policies of a subscriber.
  The [filter_policies](#smn_subscription_filter_policies) structure is documented below.

<a name="smn_subscription_filter_policies"></a>
The `filter_policies` block supports:

* `name` - (Required, String) Specifies the filter policy name. The policy name must be unique.
  + It can contain `1` to `32` characters, including lowercase letters, digits, and underscores (_).
  + It cannot start or end with an underscore, nor contain consecutive underscores. It cannot start with **smn**.

* `string_equals` - (Required, List) Specifies the string array for exact match. The array can contain `1`
  to `10` strings. The array content must be unique. The string cannot be **null** or an empty string "".
  A string can contain `1` to `32` characters, including letters, digits, and underscores (_).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The SMN subscription filter policy can be imported using `subscription_urn`, e.g.

```bash
$ terraform import huaweicloud_smn_subscription_filter_policy.test <subscription_urn>
```
