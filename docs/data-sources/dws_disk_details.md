---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_disk_details"
description: |-
  Use this data source to query DWS disk details within HuaweiCloud.
---

# huaweicloud_dws_disk_details

Use this data source to query DWS disk details within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_disk_details" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the disk details are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `instance_name` - (Optional, String) Specifies the instance name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `disk_details` - The list of disk details that matched filter parameters.  
  The [disk_details](#dws_disk_details_struct) structure is documented below.

<a name="dws_disk_details_struct"></a>
The `disk_details` block supports:

* `instance_name` - The instance name.

* `instance_id` - The instance ID.

* `host_name` - The host name.

* `disk_name` - The disk name.

* `disk_type` - The disk type.

* `total` - The total disk capacity in GB.

* `used` - The used disk capacity in GB.

* `available` - The available disk capacity in GB.

* `used_percentage` - The disk usage percentage.

* `await` - The I/O wait time in milliseconds.

* `svctm` - The I/O service time in milliseconds.

* `util` - The I/O usage percentage.

* `read_rate` - The disk read rate in KB/s.

* `write_rate` - The disk write rate in KB/s.
