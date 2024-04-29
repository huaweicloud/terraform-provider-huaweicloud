---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_event"
description: ""
---

# huaweicloud_fgs_function_event

Manages an event for testing specified function within HuaweiCloud.

## Example Usage

### Create a simple event

```hcl
variable "function_urn" {}
variable "event_name" {}
variable "event_content" {}

resource "huaweicloud_fgs_function_event" "test" {
  function_urn = var.function_urn
  name         = var.event_name
  content      = base64encode(var.event_content)
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region where the function event is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `function_urn` - (Required, String, ForceNew) Specifies the URN of the function to which the event blongs.  
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the function event name.  
  The name can contain a maximum of `25` characters and must start with a letter and end with a letter or digit.
  Only letters, digits, underscores (_) and hyphens (-) are allowed.

* `content` - (Required, String) Specifies the function event content.  
  The value is the base64 encoding of the JSON string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `updated_at` - The latest update (UTC) time of the function event, in RFC3339 format.

  -> Only events that have changed will return this attribute.

## Import

Function event can be imported using the `function_urn` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_fgs_function_event.test <function_urn>/<name>
```
