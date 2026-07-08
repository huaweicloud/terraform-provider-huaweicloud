---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_redistribution_parameters_modify"
description: |-
  Manages the modification of GaussDB redistribution parameters within HuaweiCloud.
---

# huaweicloud_gaussdb_redistribution_parameters_modify

Manages the modification of GaussDB redistribution parameters within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_redistribution_parameters_modify" "test" {
  instance_id          = var.instance_id
  redis_parallel_jobs  = 4
  redis_resource_level = "l"
}
```

### With join tables

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_redistribution_parameters_modify" "test" {
  instance_id = var.instance_id

  redis_join_tables = [
    ["database1", "schema1", "table1", "schema2", "table2"],
    ["database2", "schema3", "table3", "schema4", "table4"],
  ]
  redis_parallel_jobs  = 0
  redis_resource_level = "l"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the GaussDB instance ID.

* `redis_join_tables` - (Optional, List, NonUpdatable) Specifies the tables with JOIN relationship to enable
  multi-table expansion mode. Each element is a list of strings in the format of
  `["database", "schema1", "table1", "schema2", "table2", ...]`. Table names with special characters should be
  enclosed in double quotes. This configuration will be automatically cleared after the current expansion is completed.

* `redis_parallel_jobs` - (Optional, Int, NonUpdatable) Specifies the number of concurrent redistribution jobs.

* `redis_resource_level` - (Optional, String, NonUpdatable) Specifies the redistribution resource control level.

## Attribute Reference

The following attributes are exported:

* `id` - The resource ID, which is the instance ID.
