---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_parameters_modify"
description: |-
  Manages a TaurusDB HTAP StarRocks parameters modify resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_parameters_modify

Manages a TaurusDB HTAP StarRocks parameters modify resource within HuaweiCloud.

This is a one-time action resource that modifies parameters of a StarRocks instance by node type.
If the API response indicates that a restart is required, the resource will automatically restart
the StarRocks instance after the parameter modification job is completed.

## Example Usage

### Modify backend node parameters

```hcl
variable "taurusdb_instance_id" {}
variable "starrocks_instance_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_parameters_modify" "test" {
  taurusdb_instance_id  = var.taurusdb_instance_id
  starrocks_instance_id = var.starrocks_instance_id
  node_type             = "be"

  parameter_values = {
    "alter_tablet_worker_count"            = "10"
    "base_compaction_num_threads_per_disk" = "5"
  }
}
```

### Modify frontend node parameters

```hcl
variable "taurusdb_instance_id" {}
variable "starrocks_instance_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_parameters_modify" "test" {
  taurusdb_instance_id  = var.taurusdb_instance_id
  starrocks_instance_id = var.starrocks_instance_id
  node_type             = "fe"

  parameter_values = {
    "alter_table_timeout_second"     = "259200"
    "bdbje_heartbeat_timeout_second" = "100"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to modify the HTAP StarRocks parameters.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `taurusdb_instance_id` - (Required, String, NoneUpdatable) Specifies the TaurusDB instance ID.

* `starrocks_instance_id` - (Required, String, NoneUpdatable) Specifies the HTAP StarRocks instance ID.

* `node_type` - (Required, String, NoneUpdatable) Specifies the node type.
  The valid values are as follows:
  + **be**: backend nodes
  + **fe**: frontend nodes

* `parameter_values` - (Required, Map, NoneUpdatable) Specifies the mapping between parameter names and parameter values.
  You can specify parameter values based on a default parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
