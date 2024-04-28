---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_snapshot"
description: ""
---

# huaweicloud_dws_snapshot

Manages a GaussDB(DWS) snapshot resource within HuaweiCloud.  

## Example Usage

```hcl
  variable "cluster_id" {}
  
  resource "huaweicloud_dws_snapshot" "test" {
    name        = "demo"
    cluster_id  = var.cluster_id
    description = "This is a demo"
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Snapshot name.

  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) ID of the cluster for which you want to create a snapshot.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Snapshot description.

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
$ terraform import huaweicloud_dws_snapshot.test e87192d9-b592-4658-b23f-bdc0bb69ec2c
```
