---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_queue"
description: |-
  Manages a GaussDB(DWS) workload queue resource within HuaweiCloud.
---

# huaweicloud_dws_workload_queue

Manages a GaussDB(DWS) workload queue resource within HuaweiCloud.

## Example Usage

### Create a workload queue with CPU exclusive quotas

```hcl
variable "cluster_id" {}
variable "queue_name" {}

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id = var.cluster_id
  name       = var.queue_name

  configuration {
    resource_name  = "cpu_limit"
    resource_value = 10
  }
  configuration {
    resource_name  = "memory"
    resource_value = 0
  }
  configuration {
    resource_name  = "tablespace"
    resource_value = -1
  }
  configuration {
    resource_name  = "activestatements"
    resource_value = -1
  }
}
```

### Create a workload queue with CPU shared quotas

```hcl
variable "cluster_id" {}
variable "queue_name" {}

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id = var.cluster_id
  name       = var.queue_name

  configuration {
    resource_name  = "cpu_share"
    resource_value = 10
  }
  configuration {
    resource_name  = "memory"
    resource_value = 0
  }
  configuration {
    resource_name  = "tablespace"
    resource_value = -1
  }
  configuration {
    resource_name  = "activestatements"
    resource_value = -1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID of to which the workload queue belongs.
  Changing this parameter will create a new resource.

-> Currently, only regular cluster is supported, and logical cluster is temporarily not supported.

* `name` - (Required, String, ForceNew) Specifies the name of the workload queue, which must be unique and contains
  `3` to `28` characters, composed only of lowercase letters, numbers, or underscores (_), and must start with a
  lowercase letter. Changing this parameter will create a new resource.

* `configuration` - (Required, List, ForceNew) Specifies the configuration information for workload queue.
  Changing this parameter will create a new resource.  
  The [configuration](#DWS_workloadQueue_configuration) structure is documented below.

<a name="DWS_workloadQueue_configuration"></a>
The `configuration` block supports:

* `resource_name` - (Required, String, ForceNew) Specifies the resource name to be configured for the workload queue.  
  The valid value are as follows:
  + **memory**: memory resources.
  + **tablespace**: storage resources.
  + **activestatements**: query concurrency.
  + **cpu_limit**: exclusive quotas.
  + **cpu_share**: shared quotas.

-> When creating a workload queue, **memory**, **tablespace** and **activestatements** must be set. The **cpu_limit**
and **cpu_share** are exclusive, one of them must be set, and the **cpu_limit** is only supported for clusters above
**8.1.3**.

* `resource_value` - (Required, Int, ForceNew) Specifies the value of the resource attribute for the workload queue.
  + When the `resource name` is **memory**, the value range is from `0` to `100`, where `0` indicates no control,
    unit: %.
  + When the `resource name` is **tablespace**, the value range is from `-1` to `2,147,483,647`, where `-1` indicates
    no restriction, unit: MB.
  + When the `resource name` is **activestatements**, the value range is from `-1` to `2,147,483,647`, where `-1` and
    `0` indicates no control.
  + When the `resource name` is **cpu_limit**, the value range is from `0` to `99`, `0` means unlimited, unit: %.
  + When the `resource name` is **cpu_share**, the value range is from `1` to `99`, the default value is `20`, unit: %.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `name`.

## Import

The workload queue can be imported using `cluster_id` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_workload_queue.test <cluster_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `configuration`.
It is generally recommended running `terraform plan` after importing a workload queue.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dws_cluster" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      configuration,
    ]
  }
}
```
