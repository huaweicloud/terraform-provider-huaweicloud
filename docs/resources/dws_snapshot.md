---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_snapshot"
description: |-
  Manages a GaussDB(DWS) snapshot resource within HuaweiCloud.
---

# huaweicloud_dws_snapshot

Manages a GaussDB(DWS) snapshot resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "snapshot_name" {}
variable "description" {}

resource "huaweicloud_dws_snapshot" "test" {
  cluster_id  = var.cluster_id
  name        = var.snapshot_name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the snapshot.
  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID to which the snapshot belongs.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the snapshot.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `started_at` - Time when a snapshot starts to be created.  
  Format: ISO8601: **YYYY-MM-DDThh:mm:ssZ**.

* `finished_at` - Time when a snapshot is complete.  
  Format: ISO8601: **YYYY-MM-DDThh:mm:ssZ**.

* `size` - Snapshot size, in GB.

* `status` - Snapshot status.  
  The valid values are **CREATING**, **AVAILABLE**, and **UNAVAILABLE**.

* `type` - Snapshot type.  
  The valid values are **MANUAL**, and **AUTOMATED**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The dws snapshot can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dws_snapshot.test <id>
```
