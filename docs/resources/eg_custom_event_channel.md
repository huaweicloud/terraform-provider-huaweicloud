---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_custom_event_channel"
description: ""
---

# huaweicloud_eg_custom_event_channel

Using this resource to manage an EG custom event channel within Huaweicloud.

## Example Usage

### Manage a basic channel without enterprise project configuation

```hcl
variable "channel_name" {}

resource "huaweicloud_eg_custom_event_channel" "test" {
  name        = var.channel_name
  description = "Created by script"
}
```

### Manage a basic channel under default enterprise project

```hcl
variable "channel_name" {}

resource "huaweicloud_eg_custom_event_channel" "test" {
  name                  = var.channel_name
  description           = "Created by script"
  enterprise_project_id = "0"
}
```

### Enable cross-account configuation

```hcl
variable "channel_name" {}
variable "target_domain_ids" {
  type = list(string)
}

resource "huaweicloud_eg_custom_event_channel" "test" {
  name              = var.channel_name
  cross_account_ids = var.target_domain_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the custom event channel is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the custom event channel.  
  The valid length is limited from `1` to `128`, only letters, digits, dots (.), hyphens (-) and underscores (_) are
  allowed. The name must start with a letter or digit, and cannot be **default**.
  Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the custom event channel.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the custom
  event channel belongs.  
  The enterprise project is not used by default. Changing this will create a new resource.

* `cross_account_ids` - (Optional, List) Specifies the list of domain IDs (other tenants) for the cross-account policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The (UTC) creation time of the custom channel, in RFC3339 format.

* `updated_at` - The (UTC) update time of the custom channel, in RFC3339 format.

## Import

Custom channels can be imported by their `id` and `enterprise_project_id` (with enterprise project association), e.g.

### without enterprise project association

```bash
terraform import huaweicloud_eg_custom_event_channel.test <id>
```

### with enterprise project association

```bash
terraform import huaweicloud_eg_custom_event_channel.test <id>/<enterprise_project_id>
```
