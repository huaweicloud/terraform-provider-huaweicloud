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

* `devserver_id` - (Required, String, NonUpdatable) Specifies the ID of the DevServer.

* `action` - (Required, String, NonUpdatable) Specifies the action type of the DevServer.
  The valid values are as follows:
  + **start**: The DevServer can be started only when the DevServer is stopped, stop failure, or start failure.
  + **stop**: The DevServer can be stopped only when it is running or stop failure.
  + **reboot**: The DevServer can be rebooted only when it is running.
  + **changeos**: The DevServer OS image can be changed only when the DevServer is stopped.
  + **reinstallos**: The DevServer OS image can be reinstalled only when the DevServer is stopped.

* `admin_pass` - (Optional, String, NonUpdatable) Specifies the login password of the DevServer.

* `key_pair_name` - (Optional, String, NonUpdatable) Specifies the key pair name of the DevServer.  
  
-> Exactly one of `admin_pass` and `key_pair_name` must be set if the value of `action` parameter is **changeos** or
   **reinstallos**.

* `image_id` - (Optional, String, NonUpdatable) Specifies the image ID used to change the OS image.  
  This parameter is required when the value of `action` parameter is **changeos**.

* `user_data` - (Optional, String, NonUpdatable) Specifies the user data to be injected into the DevServer during the OS
  operation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
