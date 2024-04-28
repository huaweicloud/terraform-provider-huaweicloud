---
subcategory: "Document Database Service (DDS)"
---

# huaweicloud_dds_flavors

Use this data source to get the details of available DDS flavors.

## Example Usage

```hcl
data "huaweicloud_dds_flavors" "flavor" {
  engine_name = "DDS-Community"
  vcpus       = 8
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the flavors. If omitted,
  the provider-level region will be used.

* `engine_name` - (Required, String) Specifies the engine name. Value options: **DDS-Community** and **DDS-Enhanced**.

* `type` - (Optional, String) Specifies the type of the flavor. Value options: **mongos**, **shard**, **config**,
  **replica** and **single**.

* `vcpus` - (Optional, String) Specifies the number of vCPUs.

* `memory` - (Optional, String) Specifies the memory size in GB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `engine_name` - Indicates the engine name.

* `spec_code` - Indicates the resource specification code.

* `type` - Indicates the type of the flavor.

* `vcpus` - Indicates the number of vCPUs.

* `memory` - Indicates the memory size in GB.

* `az_status` - Indicates the mapping between availability zone and status of the flavor. **key** indicates the AZ ID,
  and **value** indicates the specification status in the AZ. Its value can be any of the following:
  + **normal**: The specification is on sale.
  + **unsupported**: This specification is not supported.
  + **sellout**: The specification is sold out.
