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

* `engine_name` - (Required, String) Specifies the engine name of the dds, "DDS-Community" and "DDS-Enhanced" are
  supported.

* `type` - (Optional, String) Specifies the type of the dds falvor. "mongos", "shard", "config", "replica" and "single"
  are supported.

* `vcpus` - (Optional, String) Specifies the vcpus of the dds flavor.

* `memory` - (Optional, String) Specifies the ram of the dds flavor in GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `spec_code` - The name of the dds flavor.
* `type` - See `type` above.
* `vcpus` - See `vcpus` above.
* `memory` - See `memory` above.
