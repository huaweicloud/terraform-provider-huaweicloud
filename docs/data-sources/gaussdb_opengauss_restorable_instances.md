---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_restorable_instances"
description: |-
  Use this data source to get the instances that can be used for backups and restorations.
---

# huaweicloud_gaussdb_opengauss_restorable_instances

Use this data source to get the instances that can be used for backups and restorations.

## Example Usage

```hcl
variable "source_instance_id" {}
variable "backup_id" {}

data "huaweicloud_gaussdb_opengauss_restorable_instances" "test" {
  source_instance_id = var.source_instance_id
  backup_id          = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `source_instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance to be restored.

* `backup_id` - (Optional, String) Specifies the instance backup ID.

* `restore_time` - (Optional, String) Specifies the time point of data restoration in the UNIX timestamp format.
  If the `backup_id` is left blank, this parameter is used to query the instance topology information and filter
  the queried instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the instances that can be used for backups and restorations.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_mode` - Indicates the instance model.
  + **enterprise**: enterprise edition
  + **standard**: standard edition
  + **basic**: basic edition

* `volume_type` - Indicates the storage type.

* `data_volume_size` - Indicates the storage space, in GB

* `version` - Indicates the instance version

* `mode` - Indicates the deployment model.
  + **Ha**: primary/standby deployment
  + **Independent**: independent deployment
