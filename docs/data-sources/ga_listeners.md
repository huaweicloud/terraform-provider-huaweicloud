---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_listeners"
description: |-
  Use this data source to get the list of listeners.
---

# huaweicloud_ga_listeners

Use this data source to get the list of listeners.

## Example Usage

```hcl
variable "listener_name" {}

data "huaweicloud_ga_listeners" "test" {
  name = var.listener_name
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Optional, String) Specifies the ID of the listener.

* `name` - (Optional, String) Specifies the name of the listener.

* `status` - (Optional, String) Specifies the current status of the listener.
  The valid values are as follows:
  + **ACTIVE**: The status of the listener is normal operation.
  + **ERROR**: The status of the listener is error.

* `accelerator_id` - (Optional, String) Specifies the ID of the accelerator to which the listener belongs.

* `protocol` - (Optional, String) Specifies the network transmission protocol type of the listener.
  The valid values are as follows:
  + **TCP**
  + **UDP**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `listeners` - The list of the listeners.
  The [listeners](#ga_listeners) structure is documented below.

<a name="ga_listeners"></a>
The `listeners` block supports:

* `id` - The ID of the listener.

* `name` - The name of the listener.  

* `description` - The description of the listener.

* `status` - The current status of the listener.

* `protocol` - The network transmission protocol type of the listener.

* `port_ranges` - The listening port range list of the listener.
  The [port_ranges](#listener_port_ranges) structure is documented below.

* `client_affinity` - The client affinity of the listener.

* `accelerator_id` - The ID of the accelerator to which the listener belongs.

* `tags` - The key/value pairs to associate with the listener.

* `created_at` - The creation time of the listener.

* `updated_at` - The latest update time of the listener.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#Listeners_frozen_info) structure is documented below.

<a name="listener_port_ranges"></a>
The `port_ranges` block supports:

* `from_port` - The listening to start port of the listener.

* `to_port` - The listening to end port of the listener.

<a name="Listeners_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.

* `effect` - The status of the resource after being forzen.

* `scene` - The service scenario.
