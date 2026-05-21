---
subcategory: "DAS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_email_batch_subscription"
description: |-
  Use this resource to batch subscribe or unsubscribe email templates within HuaweiCloud.
---

# huaweicloud_das_email_batch_subscription

Use this resource to batch subscribe or unsubscribe email templates within HuaweiCloud.

-> This resource is a one-time action resource for batch subscribing to email templates. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Subscribe to email templates

```hcl
variable "email_template_ids" {
  type = list(string)
}

resource "huaweicloud_das_email_batch_subscription" "test" {
  subscribe          = true
  email_template_ids = var.email_template_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the email templates are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `subscribe` - (Required, Bool, NonUpdatable) Specifies whether to subscribe to the email templates.  
  The valid values are as follows:
  + **true**: Subscribe to email templates.
  + **false**: Unsubscribe from email templates.

* `email_template_ids` - (Required, List, NonUpdatable) Specifies the list of email template IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
