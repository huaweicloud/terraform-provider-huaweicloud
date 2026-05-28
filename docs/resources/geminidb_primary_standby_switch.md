---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_primary_standby_switch"
description: |-
  Manages a resource to switchover GeminiDB instance node primary and standby within HuaweiCloud.
---

# huaweicloud_geminidb_primary_standby_switch

Manages a resource to switchover GeminiDB instance node primary and standby within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.
  <br/>2. The resource only supports the Primary/Standby GeminiDB Redis instance.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_primary_standby_switch" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
