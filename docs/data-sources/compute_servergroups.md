---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_servergroups"
description: ""
---

# huaweicloud_compute_servergroups

Use this data source to get the list of the compute server groups.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_compute_servergroups" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the server groups.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the server group name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servergroups` - List of ECS server groups details. The object structure of each server group is documented below.

The `servergroups` block supports:

* `id` - The server group ID in UUID format.

* `name` - The server group name.

* `policies` - The set of policies for the server group.

* `members` - An array of one or more instance ID attached to the server group.
