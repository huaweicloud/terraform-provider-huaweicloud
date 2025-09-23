---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "Huaweicloud: huaweicloud_eg_event_subscription_batch_action"
description: |-
  Use this resource to operate the EG event subscription within HuaweiCloud.
---

# huaweicloud_eg_event_subscription_batch_action

Use this resource to operate the EG event subscription within HuaweiCloud.

-> This resource is only a one-time action resource for operating event subscription status. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "subscription_ids" {
  type = list(string)
}

resource "huaweicloud_eg_event_subscription_batch_action" "example" {
  subscription_ids = var.subscription_ids
  operation        = "ENABLE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the event subscriptions are located.  
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `subscription_ids` - (Required, List, NonUpdatable) Specifies the list of subscription IDs to be operated.  
  The single operation only can handle up to `10` event subscriptions at most.

* `operation` - (Required, String, NonUpdatable) Specifies whether to enable the event subscription.
  The valid values are as follows:
  + **ENABLE**
  + **DISABLE**

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project to which the
subscriptions belong.  
  This parameter is only valid for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
