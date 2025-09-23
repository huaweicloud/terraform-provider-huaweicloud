---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_devserver_action"
description: |-
  Use this resource to operate ModelArts DevServer within HuaweiCloud.
---
# huaweicloud_modelarts_devserver_action

Use this resource to operate ModelArts DevServer within HuaweiCloud.

## Example Usage

```hcl
variable "devserver_id" {}

resource "huaweicloud_modelarts_devserver_action" "test" {
  devserver_id = var.devserver_id
  action       = "start"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `devserver_id` - (Required, String, ForceNew) Specifies the ID of the DevServer.
  Changing this creates a new resource.

* `action` - (Required, String, ForceNew) Specifies the action type of the DevServer.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **start**: The DevServer can be started only when the DevServer is stopped, stop failure, or start failure.
  + **stop**: The DevServer can be stopped only when it is running or stop failure.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
