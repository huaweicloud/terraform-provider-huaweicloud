---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_object_compare"
description: |-
  Manages a resource to start a DRS object compare task within HuaweiCloud.
---

# huaweicloud_drs_object_compare

Manages a resource to start a DRS object compare task within HuaweiCloud.

-> This resource is a one-time action resource used to start an object compare task. Deleting this resource will not
  undo the compare operation, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "job_id" {} 
variable "compare_task_num" {}

resource "huaweicloud_drs_object_compare" "test" { 
  job_id           = var.job_id 
  compare_task_num = var.compare_task_num
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the DRS job ID.

* `compare_task_num` - (Required, Int, NonUpdatable) Specifies the number of concurrent compare task threads.
  This parameter is currently effective only for cloudDataGuard-cassandra and
  cloudDataGuard-gausscassandra-to-gausscassandra links.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
