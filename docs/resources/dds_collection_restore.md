---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_collection_restore"
description: |-
  Manages a DDS collection restore resource within HuaweiCloud.
---

# huaweicloud_dds_collection_restore

Manages a DDS collection restore resource within HuaweiCloud.

## Example Usage

### Restore database level data to instance

```hcl
variable "instance_id" {}
variable "database" {}
variable "restore_database_time" {}

resource "huaweicloud_dds_collection_restore" "test" {
  instance_id  = var.instance_id

  restore_collections {
    database              = var.database
    restore_database_time = var.restore_database_time
  }
}
```

### Restore collection level data to instance

```hcl
variable "instance_id" {}
variable "database" {}
variable "collection" {}
variable "restore_collection_time" {}

resource "huaweicloud_dds_collection_restore" "test" {
  instance_id  = var.instance_id

  restore_collections {
    database = var.database

    collections {
      old_name                = var.collection
      restore_collection_time = var.restore_collection_time
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `restore_collections` - (Required, List, ForceNew) Specifies the restore informations.
  Changing this creates a new resource.
  The [restore_collections](#block--restore_collections) structure is documented below.

<a name="block--restore_collections"></a>
The `restore_collections` block supports:

* `database` - (Required, String, ForceNew) Specifies the database name.
  Changing this creates a new resource.

* `collections` - (Optional, List, ForceNew) Specifies the collection informations.
  Changing this creates a new resource.
  The [collections](#block--restore_collections--collections) structure is documented below.

* `restore_database_time` - (Optional, String, ForceNew) Specifies the database restoration time point.
  It is mandatory for database level restoration. The value is a UNIX timestamp, in milliseconds. The time zone is UTC.
  Changing this creates a new resource.

<a name="block--restore_collections--collections"></a>
The `collections` block supports:

* `old_name` - (Required, String, ForceNew) Specifies the collection name before the restoration.
  Changing this creates a new resource.

* `restore_collection_time` - (Required, String, ForceNew) Specifies the collection restoration time point.
  The value is a UNIX timestamp, in milliseconds. The time zone is UTC.
  Changing this creates a new resource.

* `new_name` - (Optional, String, ForceNew) Specifies the collection name after the restoration.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
