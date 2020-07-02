---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_flavors_v3"
sidebar_current: "docs-huaweicloud-datasource-dds-flavors-v3"
description: |-
  Get the flavor information on HuaweiCloud DDS service.
---

# huaweicloud\_dds\_flavors\_v3

Use this data source to get the ID of an available HuaweiCloud dds flavor.

## Example Usage

```hcl
data "huaweicloud_dds_flavors_v3" "flavor" {
    engine_name = "DDS-Community"
    vcpus = 8
}
```

## Argument Reference

* `region` - (Optional) Specifies the region in which to obtain the V3 dds client.

* `engine_name` - (Required) Specifies the engine name of the dds, "DDS-Community" and "DDS-Enhanced" are supported.

* `type` - (Optional) Specifies the type of the dds falvor. "mongos", "shard", "config", "replica" and "single" are supported.

* `vcpus` - (Optional) Specifies the vcpus of the dds flavor.

* `memory` - (Optional) Specifies the ram of the dds flavor in GB.


## Attributes Reference

* `region` - See Argument Reference above.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `spec_code - The name of the rds flavor.
* `type` - See `type` above.
* `vcpus` - See `vcpus` above.
* `memory` - See 'memory' above.
