---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_redis_flavors"
description: |-
  Use this data source to get list of GaussDB redis flavors.
---

# huaweicloud_gaussdb_redis_flavors

Use this data source to get list of GaussDB redis flavors.

## Example Usage

```hcl
data "huaweicloud_gaussdb_redis_flavors" "flavors" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of flavors.

  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block contains:

* `engine_name` - Indicates the database type.

* `engine_version` - Indicates the database version.

* `vcpus` - Indicates the number of vCPUs.

* `ram` - Indicates the memory size in megabytes (MB).

* `spec_code` - Indicates the resource specification code.

* `az_status` - Indicates the status of specifications in an AZ. The value can be:
  + **normal**: indicating that the specifications are on sale.
  + **unsupported**: indicating that the specifications are not supported.
  + **sellout**: indicating that the specifications are sold out.
