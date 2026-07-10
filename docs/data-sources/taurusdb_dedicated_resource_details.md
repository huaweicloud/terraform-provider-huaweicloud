---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_dedicated_resource_details"
description: |-
  Use this data source to query the detail of a TaurusDB dedicated resource within HuaweiCloud.
---

# huaweicloud_taurusdb_dedicated_resource_details

Use this data source to query the detail of a TaurusDB dedicated resource within HuaweiCloud.

## Example Usage

```hcl
variable "dedicated_resource_id" {}

data "huaweicloud_taurusdb_dedicated_resource_details" "test" {
  dedicated_resource_id = var.dedicated_resource_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region
  will be used.

* `dedicated_resource_id` - (Required, String) Specifies the ID of the dedicated resource pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_name` - Indicates the name of the dedicated resource pool.

* `engine_name` - Indicates the DB engine name.

* `availability_zone_ids` - Indicates the list of AZs.

* `architecture` - Indicates the CPU architecture.

* `status` - Indicates the status of the dedicated resource pool. The value can be **NORMAL**, **BUILDING** or **DELETED**.

* `dedicated_compute_info` - Indicates the compute resource information.
  The [dedicated_compute_info](#dedicated_compute_info_struct) structure is documented below.

* `dedicated_storage_info` - Indicates the storage resource information.
  The [dedicated_storage_info](#dedicated_storage_info_struct) structure is documented below.

<a name="dedicated_compute_info_struct"></a>
The `dedicated_compute_info` block contains:

* `vcpus_total` - Indicates the total vCPUs in the dedicated resource pool.

* `vcpus_used` - Indicates the used vCPUs in the dedicated resource pool.

* `ram_total` - Indicates the total memory size of the dedicated resource pool, in GB.

* `ram_used` - Indicates the used memory size of the dedicated resource pool, in GB.

* `spec_code` - Indicates the compute resource specification code of the dedicated resource pool.

* `host_num` - Indicates the number of compute hosts in the dedicated resource pool.

<a name="dedicated_storage_info_struct"></a>
The `dedicated_storage_info` block contains:

* `spec_code` - Indicates the storage resource specification code of the dedicated resource pool.

* `host_num` - Indicates the number of storage hosts in the dedicated resource pool.
