---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_eventrouter_cluster"
description: "Manages an EG event router cluster within HuaweiCloud."
---

# huaweicloud_eg_eventrouter_cluster

Using this resource to manage an EG event router cluster within HuaweiCloud.

## Example Usage

### Create an event router cluster with basic configuration

```hcl
variable "cluster_name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "availability_zones" {
  type = list(string)
}

resource "huaweicloud_eg_eventrouter_cluster" "test" {
  name               = var.cluster_name
  source_type        = "KAFKA"
  sink_type          = "KAFKA"
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  description        = "Created by terraform script"
  availability_zones = join(",", var.availability_zones)
  flavor             = "small"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the event router cluster is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the event router cluster.  
  The name can contain maximum of 64 characters, which may consist of letters, digits, underscores (_), and hyphens (-).

* `source_type` - (Required, String, NonUpdatable) Specifies the source type of the event router cluster.  
  The valid values are as follows:
  + **KAFKA** - Kafka event source
  + **DMS_ROCKETMQ** - RocketMQ event source
  + **DMS_RABBITMQ** - RabbitMQ event source

* `sink_type` - (Required, String, NonUpdatable) Specifies the sink type of the event router cluster.  
  The valid values are as follows:
  + **KAFKA** - Kafka event sink
  + **DMS_ROCKETMQ** - RocketMQ event sink
  + **DMS_RABBITMQ** - RabbitMQ event sink

  -> The value of `sink_type` parameter must be same as the value of `source_type` parameter.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID to which the event router cluster belongs.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID to which the event router cluster belongs.

* `description` - (Optional, String) Specifies the description of the event router cluster.  
  The description can contain a maximum of `255` characters.

* `availability_zones` - (Optional, String, NonUpdatable) Specifies the availability zone names of the event router
  cluster.  
  The multiple availability zones are separated by commas (,), e.g. **cn-north-4a,cn-north-4g**.

* `flavor` - (Optional, String, NonUpdatable) Specifies the flavor of the event router cluster.  
  The valid values are as follows:
  + **small** - Small flavor (`1` vCPU, `2` GB memory)
  + **medium** - Medium flavor (`2` vCPU, `4` GB memory)
  + **large** - Large flavor (`4` vCPU, `8` GB memory)
  + **xlarge** - Extra large flavor (`8` vCPU, `16` GB memory)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `status` - The status of the event router cluster.
  + **RUNNING** - The cluster is running normally
  + **ERROR** - The cluster creation failed

* `job_count` - The number of jobs running in the event router cluster.

* `created_at` - The creation time of the event router cluster, in RFC3339 format.

* `updated_at` - The latest update time of the event router cluster, in RFC3339 format.

## Import

Event router clusters can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_eg_eventrouter_cluster.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `flavor`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition should be updated to
align with the cluster.

```hcl
resource "huaweicloud_eg_eventrouter_cluster" "test" {
  ...

  lifecycle {
    ignore_changes = [
      flavor,
    ]
  }
}
```
