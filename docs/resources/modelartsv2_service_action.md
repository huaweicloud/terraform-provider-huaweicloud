---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_service_action"
description: |-
  Use this resource to operate the ModelArts service within HuaweiCloud.
---

# huaweicloud_modelartsv2_service_action

Use this resource to operate the ModelArts service within HuaweiCloud.

-> This resource is only a one-time action resource for operating the service. Deleting this resource will not undo
   action that has been performed, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "service_id" {}

resource "huaweicloud_modelartsv2_service_action" "test" {
  service_id = var.service_id
  action     = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the service is located.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `service_id` - (Required, String, NonUpdatable) Specifies the service ID to be operated.

* `action` - (Required, String, NonUpdatable) Specifies the action type.  
  The valid values are as follows:
  + **start**
  + **stop**
  + **interrupt**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
