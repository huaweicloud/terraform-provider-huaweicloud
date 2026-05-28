---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_available_flavors"
description: |-
  Use this data source to get the list of GeminiDB instance available flavors.
---

# huaweicloud_geminidb_available_flavors

Use this data source to get the list of GeminiDB instance available flavors.

-> This data source supports GeminiDB Cassandra, GeminiDB Mongo, GeminiDB Influx and GeminiDB Redis instances.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_available_flavors" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `current_flavor` - The instance current specification information.
  The [current_flavor](#current_flavor_struct) structure is documented below.

* `optional_flavors` - The list of available specification options that the instance specifications can be changed to.
  The [optional_flavors](#optional_flavors_struct) structure is documented below.

<a name="current_flavor_struct"></a>
The `current_flavor` block supports:

* `vcpus` - The number of CPUs.

* `ram` - The memory size, in GB.

* `spec_code` - The specification code.

* `az_status` - The AZ status.
  + **normal**: Available.
  + **abandon**: Offline.
  + **sellout**: Sold-out.

<a name="optional_flavors_struct"></a>
The `optional_flavors` block supports:

* `vcpus` - The number of CPUs.

* `ram` - The memory size, in GB.

* `spec_code` - The specification code.

* `az_status` - The AZ status.
  + **normal**: Available.
  + **abandon**: Offline.
  + **sellout**: Sold-out.
