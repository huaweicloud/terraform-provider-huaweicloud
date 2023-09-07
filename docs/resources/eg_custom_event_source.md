---
subcategory: "EventGrid (EG)"
---

# huaweicloud_eg_custom_event_source

Using this resource to manage an EG custom event source within Huaweicloud.

## Example Usage

```hcl
variable "channel_id" {}
variable "source_name" {}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = var.channel_id
  type        = "RABBITMQ"
  name        = var.source_name
  description = "Created by script"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the custom event channel and custom event source
  are located. If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `channel_id` - (Required, String, ForceNew) Specifies the ID of the custom event channel to which the custom event
  source belongs.  
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the custom event source.  
  The valid length is limited from `1` to `128`, only lowercase letters, digits, hyphens (-), underscores (_) are
  allowed. The name must start with a lowercase letter or digit.  
  Changing this will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the custom event source.
  The valid values are as follows:
  + **APPLICATION**
  + **RABBITMQ**
  + **ROCKETMQ**

  Defaults to **APPLICATION**.  
  Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the custom event source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - The status of the custom event source.

* `created_at` - The (UTC) creation time of the custom event source, in RFC3339 format.

* `updated_at` - The (UTC) update time of the custom event source, in RFC3339 format.

## Import

Custom event sources can be imported by their `id`, e.g.

```bash
terraform import huaweicloud_eg_custom_event_source.test <id>
```
