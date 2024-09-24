---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_snapshot_policy"
description: |-
  Manages a GaussDB(DWS) automated snapshot policy resource within HuaweiCloud.  
---

# huaweicloud_dws_snapshot_policy

Manages a GaussDB(DWS) automated snapshot policy resource within HuaweiCloud.  

## Example Usage

### Create a policy to periodically create a full snapshot

```hcl
variable "cluster_id" {}

resource "huaweicloud_dws_snapshot_policy" "test" {
  name       = "demo"
  cluster_id = var.cluster_id
  type       = "full"
  strategy   = "0 8 6 4 * ?"
}
```

### Create a policy to periodically create a full snapshot at specified time

```hcl
  variable "cluster_id" {}

  resource "huaweicloud_dws_snapshot_policy" "test" {
    name       = "demo"
    cluster_id = var.cluster_id
    type       = "full"
    strategy   = "2023-05-19T09:24:00"
  }
```

### Create a policy to periodically create a increment snapshot

```hcl
  variable "cluster_id" {}

  resource "huaweicloud_dws_snapshot_policy" "test" {
    name       = "demo"
    cluster_id = var.cluster_id
    type       = "increment"
    strategy   = "0 8 6 4 * ?"
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) The cluster ID of which the automated snapshot policy belongs to.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of the automated snapshot policy.

  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) The type of the automated snapshot policy.  
  The options are as follows:
    + **full**: A full snapshot backs up the data of an entire cluster.
    + **increment**: An incremental snapshot records the changes made after the previous snapshot was created.

  Changing this parameter will create a new resource.

* `strategy` - (Required, String, ForceNew) The strategy of the automated snapshot policy.  
  Its format is a Cron expression, which specifies when to create a snapshot.

  Changing this parameter will create a new resource.

  -> Note: The UTC time is used by default. Set the policy based on the time zone and time difference as required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The DWS snapshot policy can be imported using `cluster_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_snapshot_policy.test <cluster_id>/<id>
```
