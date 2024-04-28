---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network"
description: ""
---

# huaweicloud_cc_central_network

Manages a central network resource of Cloud Connect within HuaweiCloud.

## Example Usage

```hcl
variable "central_network_name" {}

resource "huaweicloud_cc_central_network" "test" {
  name        = var.central_network_name
  description = "This is a demo"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the central network.

* `description` - (Optional, String) The description of the central network.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project ID to which the central network belongs.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) The key/value pairs to associate with the central network.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `state` - The state of the central network.
  The valid values are as follows:
    - AVAILABLE
    - UPDATING
    - FAILED
    - CREATING
    - DELETING
    - DELETED

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The central network can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_central_network.test 0ce123456a00f2591fabc00385ff1234
```
