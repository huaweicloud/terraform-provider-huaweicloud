---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_support_links"
description: |-
  Use this data source to get a list of available support links for DRS.
---

# huaweicloud_drs_support_links

Use this data source to get a list of available support links for DRS.

## Example Usage

```hcl
variable "job_type" {}

data "huaweicloud_drs_support_links" "test" {
  job_type = var.job_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_type` - (Required, String) Specifies the DRS job type.
  The valid values are as follows:
  + **migration**: Online migration.
  + **sync**: Data synchronization.
  + **cloudDataGuard**: Disaster recovery.
  + **replay**: Recording and playback.
  + **verify**: Verification task.
  + **cdc**: CDC task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `support_links` - The list of available support links.

The [support_links](#support_links_struct) structure is documented below.

<a name="support_links_struct"></a>
The `support_links` block supports:

* `engine_type` - The engine type.

* `net_type` - The network type.
  The valid values are as follows:
  + **eip**: Public network.
  + **vpc**: VPC network.
  + **vpn**: VPN and dedicated network.

* `task_modes` - The list of task modes.
  The valid values are as follows:
  + **FULL_TRANS**: Full migration.
  + **FULL_INCR_TRANS**: Full + Incremental migration.
  + **INCR_TRANS**: Incremental migration.

* `job_direction` - The direction of data flow.
  The valid values are as follows:
  + **up**: In cloud, corresponds to standby in disaster recovery scenario.
  + **down**: Out cloud, corresponds to primary in disaster recovery scenario.
  + **non-dbs**: Self-built database.

* `cluster_mode` - The cloud instance type.
  The valid values are as follows:
  + **Single**: Single machine mode.
  + **Ha**: Primary/standby mode.
  + **Cluster**: Cluster mode.
  + **Sharding**: Sharding mode.
  + **Independent**: GaussDB independent deployment mode.

* `job_instance_type` - The DRS instance type.
  The valid values are as follows:
  + **Single**: Single instance.
  + **Ha**: Primary/standby instance.

* `is_support_bind_eip` - Whether binding EIP is supported.
