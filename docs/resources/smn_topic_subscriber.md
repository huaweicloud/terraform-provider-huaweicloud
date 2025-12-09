---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topic_subscriber"
description: |-
  Manages an SMN topic subscriber resource within HuaweiCloud.
---

# huaweicloud_smn_topic_subscriber

Manages an SMN topic subscriber resource within HuaweiCloud.

## Example Usage

```hcl
variable "topic_urn" {}
variable "subscribe_id" {}

resource "huaweicloud_smn_topic_subscriber" "test" {
  topic_urn    = var.topic_urn
  subscribe_id = var.subscribe_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN topic subscriber resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the unique resource identifier of the topic.

* `subscribe_id` - (Required, String, NonUpdatable) Specifies the ID of the subscriber.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `protocol` - Indicates the subscription protocol.

* `owner` - Indicates the project ID of the topic creator.

* `endpoint` - Indicates the message receiving endpoint.

* `remark` - Indicates the remarks.

* `status` - Indicates the subscription status.
  + **0**: indicates that the subscription is not confirmed.
  + **1**: indicates that the subscription is confirmed.
  + **3**: indicates that the subscription is canceled.

* `filter_policies` - Indicates the filter policy.
  The [filter_policies](#filter_policies_struct) structure is documented below.

<a name="filter_policies_struct"></a>
The `filter_policies` block supports:

* `name` - Indicates the filter policy name.

* `string_equals` - Indicates the string array for exact match.

## Import

This resource can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_smn_topic_subscriber.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `subscribe_id`. It is generally recommended
running `terraform plan` after importing a topic subscriber. You can then decide if changes should be applied to the
topic subscriber, or the resource definition should be updated to align with the topic subscriber. Also you can ignore
changes as below.

```hcl
resource "huaweicloud_smn_topic_subscriber" "test" {
  ...

  lifecycle {
    ignore_changes = [
      subscribe_id,
    ]
  }
}
```
