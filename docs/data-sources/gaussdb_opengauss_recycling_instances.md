---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_recycling_instances"
description: |-
  Use this data source to get the list of GaussDB OpenGauss instances in the recycle bin.
---

# huaweicloud_gaussdb_opengauss_recycling_instances

Use this data source to get the list of GaussDB OpenGauss instances in the recycle bin.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_recycling_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_name` - (Optional, String) Specifies the GaussDB OpenGauss instance name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the information about all instances in the recycle bin.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `mode` - Indicates the product type.
  The value can be:
  + **basic**: basic edition
  + **standard**: standard edition
  + **enterprise**: enterprise edition

* `ha_mode` - Indicates the deployment model.
  The value can be:
  + **Ha**: primary/standby deployment
  + **Independent**: independent deployment
  + **Combined**: combined deployment

* `engine_name` - Indicates the engine name.

* `engine_version` - Indicates the engine version.

* `pay_model` - Indicates the billing mode.
  The value can be:
  + **0**: pay-per-use
  + **1**: yearly/monthly

* `volume_type` - Indicates the disk type.
  The value can be:
  + **high**: high I/O
  + **ultrahigh**: ultra-high I/O
  + **essd**: extreme SSD

* `volume_size` - Indicates the disk size.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `enterprise_project_name` - Indicates the enterprise project name.

* `recycle_backup_id` - Indicates the backup ID.

* `backup_level` - Indicates the backup level.

* `data_vip` - Indicates the private IP address.

* `recycle_status` - Indicates the backup status in the recycle bin.
  The value can be: **Running**, **Active**.

* `created_at` - Indicates the creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `deleted_at` - Indicates the deletion time in the **yyyy-mm-ddThh:mm:ssZ** format.
