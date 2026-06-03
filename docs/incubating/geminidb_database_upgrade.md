---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_database_upgrade"
description: |-
  Manages a resource to upgrade database patch within HuaweiCloud.
---

# huaweicloud_geminidb_database_upgrade

Manages a resource to upgrade database patch within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request
  record, but will only remove the resource information from the tf state file.
  <br/>2. The instance types that support the database patch upgrade are as follows:
  **GeminiDB Cassandra**, **GeminiDB Influx** and **GeminiDB Redis**.
  <br/>3. The resource is not available to frozen or abnormal instances.
  <br/>4. The resource is not available if there are abnormal instance nodes.
  <br/>5. You can view field `patch_available` in the result returned by the API for querying instance details and
  check whether the database patch upgrade is supported.
  <br/>6. Perform an upgrade during off-peak hours.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_database_upgrade" "test" {
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

* `create` - Default is 60 minutes.
