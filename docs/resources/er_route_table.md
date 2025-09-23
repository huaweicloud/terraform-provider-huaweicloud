---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_route_table"
description: ""
---

# huaweicloud_er_route_table

Manages a route table resource under the ER instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "route_table_name" {}

resource "huaweicloud_er_route_table" "test" {
  instance_id = var.instance_id
  name        = var.route_table_name
  description = "Route table created by terraform"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ER instance and route table are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the ER instance to which the route table belongs.  
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the route table.  
  The name can contain `1` to `64` characters, only English letters, Chinese characters, digits, underscore (_),
  hyphens (-) and dots (.) allowed.

* `description` - (Optional, String) Specifies the description of the route table.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the route table.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `is_default_association` - Whether this route table is the default association route table.

* `is_default_propagation` - Whether this route table is the default propagation route table.

* `status` - The current status of the route table.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Route tables can be imported using their `id` and the related `instance_id`, separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_er_route_table.test <instance_id>/<id>
```
