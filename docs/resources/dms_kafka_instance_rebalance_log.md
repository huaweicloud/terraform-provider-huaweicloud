---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_rebalance_log"
description: |-
  Manages a Kafka rebalance log resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_rebalance_log

Manages a Kafka rebalance log resource within HuaweiCloud.

-> 1. Only one `huaweicloud_dms_kafka_instance_rebalance_log` resource can be created under a Kafka instance.
   <br>2. Single-node instance do not support rebalance log.
   <br>3. Deleting the resource only disables rebalance logging. The existing LTS log groups and log streams will
   remain and will still incur charges.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_kafka_instance_rebalance_log" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Kafka instance rebalance log is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance to which the
  rebalance log belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the `instance_id`.

* `log_group_id` - The ID of the log group.

* `log_stream_id` - The ID of the log stream.

* `dashboard_id` - The ID of the dashboard.

* `log_type` - The type of the log.

* `log_file_name` - The name of the log file.

* `status` - The status of the rebalance log.

* `created_at` - The creation time of the rebalance log, in RFC3339 format.

* `updated_at` - The latest update time of the rebalance log, in RFC3339 format.

## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

This resource can be imported using its `id`, e.g.

```bash
terraform import huaweicloud_dms_kafka_instance_rebalance_log.test <id>
```
