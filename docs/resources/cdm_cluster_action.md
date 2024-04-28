---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_cluster_action"
description: ""
---

# huaweicloud_cdm_cluster_action

Manages a CDM cluster action resource within HuaweiCloud.  

## Example Usage

### Restart service process immediately

```hcl
variable "cdm_cluster_id" {}

resource "huaweicloud_cdm_cluster_action" "restart" {
  cluster_id = var.cdm_cluster_id
  type       = "restart"

  restart {
    level      = "SERVICE"
    mode       = "IMMEDIATELY"
  }
}
```

### Restart vm

```hcl
variable "cdm_cluster_id" {}

resource "huaweicloud_cdm_cluster_action" "restart" {
  cluster_id = var.cdm_cluster_id
  type       = "restart"

  restart {
    level      = "VM"
    mode       = "FORCIBLY"
    delay_time = 0
  }
}
```

### Start a cluster

```hcl
variable "cdm_cluster_id" {}

resource "huaweicloud_cdm_cluster_action" "start" {
  cluster_id = var.cdm_cluster_id
  type       = "start"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) ID of CDM cluster.  

  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Action type.  
  Value options are as follows:
    + **start**: start the cluster.
    + **restart**: restart the service process or VMs of cluster.

  Changing this parameter will create a new resource.

* `restart` - (Optional, List, ForceNew) The configuration of the restart action.  
  This field is required when the type is **restart**.

  Changing this parameter will create a new resource.
The [restart](#CdmClusterAction_Restart) structure is documented below.

<a name="CdmClusterAction_Restart"></a>
The `restart` block supports:

* `level` - (Required, String, ForceNew) Restart level.  
  Value options are as follows:
    + **SERVICE**: service restart.
    + **VM**: VM restart.

  Changing this parameter will create a new resource.

* `mode` - (Required, String, ForceNew) Restart mode.  
  Value options are as follows:
    + **IMMEDIATELY**: immediate restart.
    + **FORCIBLY**: forcible restart.
      Restarte the service process will interrupt the service process and restart the VMs in the cluster.
    + **SOFTLY**: common restart.

  Changing this parameter will create a new resource.

* `delay_time` - (Optional, Int, ForceNew) Restart delay, in seconds.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
