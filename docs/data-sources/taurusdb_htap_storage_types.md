---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_storage_types"
description: |-
  Use this data source to query the list of TaurusDB HTAP instance storage types within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_storage_types

Use this data source to query the list of TaurusDB HTAP instance storage types within HuaweiCloud.

## Example Usage

```hcl
variable "version_name" {}

data "huaweicloud_taurusdb_htap_storage_types" "test" {
  database     = "star-rocks"
  version_name = var.version_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP instance storage types.
  If omitted, the provider-level region will be used.

* `database` - (Required, String) Specifies the HTAP data engine type. Value options: **star-rocks**.

* `version_name` - (Required, String) Specifies the major version of the HTAP database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_type` - The list of HTAP instance storage type that matched filter parameters.
  The [storage_type](#taurusdb_htap_storage_type_attr) structure is documented below.

<a name="taurusdb_htap_storage_type_attr"></a>
The `storage_type` block supports:

* `name` - The name of the storage medium.

* `az_status` - Map of availability zone names and their status for the storage type.
  The valid values are as follows:
  + **normal**: The storage type is available in the AZ.
  + **unsupported**: The storage type is not supported.
  + **sellout**: The storage type is sold out.

* `min_volume_size` - The minimum disk size in GB for the storage type.
