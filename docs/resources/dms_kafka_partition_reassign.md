---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_partition_reassign"
description: |-
  Manages a DMS kafka partition reassign resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_partition_reassign

Manages a DMS kafka partition reassign resource within HuaweiCloud.

## Example Usage

### Create a partition reassignment task by manually specified assignment plan

```hcl
variable "instance_id" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_partition_reassign" "test" {
  instance_id = var.instance_id
  
  reassignments {
    topic = var.topic_name

    assignment {
      partition         = 0
      partition_brokers = [0,1,2]
    }

    assignment {
      partition         = 1
      partition_brokers = [2,0,1]
    }

    assignment {
      partition         = 2
      partition_brokers = [1,2,0]
    }
  }
}
```

### Create a partition reassignment task by automatic assignment plan

```hcl
variable "instance_id" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_partition_reassign" "test" {
  instance_id = var.instance_id
  
  reassignments {
    topic              = var.topic_name
    brokers            = [0,1,2]
    replication_factor = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DMS kafka instance ID.
  Changing this creates a new resource.

* `reassignments` - (Required, List, ForceNew) Specifies the reassignment plan.
  Changing this creates a new resource.
  The [reassignments](#reassignments_struct) structure is documented below.

* `throttle` - (Optional, Int, ForceNew) Specifies the reassignment threshold. Value can be specified ranges from `1`
  to `300`. The unit is **MB/s**. Or specifies it to `-1`, indicating no throttling required.
  Changing this creates a new resource.

* `is_schedule` - (Optional, Bool, ForceNew) Specifies whether the task is scheduled. Defaults to **false**.
  Changing this creates a new resource.

* `execute_at` - (Optional, Int, ForceNew) Specifies the schedule time. The value is a UNIX timestamp, in **ms**.
  It's required if `is_schedule` is **true**. Changing this creates a new resource.

* `time_estimate` - (Optional, Bool, ForceNew) Specifies whether to perform time estimation tasks. Defaults to **false**.
  Changing this creates a new resource.

<a name="reassignments_struct"></a>
The `reassignments` block supports:

* `topic` - (Required, String, ForceNew) Specifies the topic name. Changing this creates a new resource.

* `brokers` - (Optional, List, ForceNew) Specifies the integer list of brokers to which partitions are reassigned.
  It's **required** in **automatic** assignment. Changing this creates a new resource.

* `replication_factor` - (Optional, Int, ForceNew) Specifies the replication factor, which can be specified in
  **automatic** assignment. Changing this creates a new resource.

* `assignment` - (Optional, List, ForceNew) Specifies the manually specified assignment plan.
  It's **required** in **manually** specified assignment. Changing this creates a new resource.
  The [assignment](#reassignments_assignment_struct) structure is documented below.

-> If manually specified assignment and automatic assignment are both specified, only **manually** specified assignment
will take effect.

<a name="reassignments_assignment_struct"></a>
The `assignment` block supports:

* `partition` - (Optional, Int, ForceNew) Specifies the partition number in manual assignment.
  It's actually **required** in **manual** assignment plan. Changing this creates a new resource.

* `partition_brokers` - (Optional, List, ForceNew) Specifies the integer list of brokers to be assigned to a partition in
  manual assignment. It's actually **required** in **manual** assignment plan. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `task_id` - Indicates the task ID, and it's only returned for a partition reassignment task.

* `reassignment_time` - Indicates the estimated time, in seconds, and it's only returned for a time estimation task.

## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 20 minutes.
