---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow_subscription"
description: |-
  Manages a ModelArts workflow message subscription configuration within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow_subscription

Manages a ModelArts workflow message subscription configuration within HuaweiCloud.

## Example Usage

```hcl
variable "workflow_id" {}
variable "subscribe_topic_urns" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_workflow_subscription" "test" {
  workflow_id = var.workflow_id
  topic_urns  = var.subscribe_topic_urns

  events = [
    "service_step:wait_inputs,hold,completed",
    "labeling:wait_inputs,hold,failed,create_failed",
    "*:wait_inputs,hold,completed",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workflow subscription is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workflow_id` - (Required, String, NonUpdatable) Specifies the ID of the workflow to which the subscription belongs.

* `topic_urns` - (Required, List) Specifies the list of SMN topic URNs to subscribe.

* `events` - (Optional, List) Specifies the list of workflow events to subscribe.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (subscription ID).

* `created_at` - The creation time of the workflow subscription, in RFC3339 format.

## Import

Workflow subscriptions can be imported using `workflow_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_modelartsv2_workflow_subscription.test <workflow_id>/<id>
```
