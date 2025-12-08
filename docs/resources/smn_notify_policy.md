---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_notify_policy"
description: |-
  Manages an SMN notify policy resource within HuaweiCloud.
---

# huaweicloud_smn_notify_policy

Manages an SMN notify policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "topic_urn" {}
variable "subscription_urn" {}

resource "huaweicloud_smn_notify_policy" "test" {
  topic_urn = var.topic_urn
  protocol = "callnotify"
  polling {
    order             = 1
    subscription_urns = [var.subscription_urn]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN notify policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the resource identifier of the topic.

* `protocol` - (Required, String) Specifies the notification policy type. Only the voice notification policies are
  supported. Value options: **callnotify**.

* `polling` - (Required, List) Specifies the subscription endpoint in a polling notification policy.
  The [polling](#polling_struct) structure is documented below.

<a name="polling_struct"></a>
The `polling` block supports:

* `order` - (Required, Int) Specifies the sequence number of the subscription endpoint being polled.

* `subscription_urns` - (Required, List) Specifies the URN list of subscription endpoints.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `polling` - Indicates the subscription endpoint in a polling notification policy.
  The [polling](#polling_attribute) structure is documented below.

<a name="polling_attribute"></a>
The `polling` block supports:

* `subscriptions` - Indicates the URN list of subscription endpoints.
  The [subscriptions](#subscriptions_attribute) structure is documented below.

<a name="subscriptions_attribute"></a>
The `subscriptions` block supports:

* `subscription_urn` - Indicates the sequence number of the subscription endpoint being polled.

* `endpoint` - Indicates the URN list of subscription endpoints.

* `remark` - Indicates the remark.

* `status` - Indicates the subscription status. The value can be:
  + **0**: indicates the subscription has not been confirmed.
  + **1**: indicates that the subscription has been confirmed.
  + **3**: indicates that the subscription has been canceled.

## Import

This resource can be imported using the `topic_urn`, e.g.:

```bash
$ terraform import huaweicloud_smn_notify_policy.test <topic_urn>
```
