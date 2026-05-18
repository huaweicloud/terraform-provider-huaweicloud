---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_dedicated_resources"
description: |-
  Use this data source to get the list of dedicated resource.
---

# huaweicloud_geminidb_dedicated_resources

Use this data source to get the list of dedicated resource.

-> This data source only supports GeminiDB Cassandra instances.

## Example Usage

```hcl
data "huaweicloud_geminidb_dedicated_resources" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of dedicated resources.
  The [resources](#dedicated_resources_struct) structure is documented below.

<a name="dedicated_resources_struct"></a>
The `resources` block supports:

* `id` - The dedicated resource ID.

* `resource_name` - The dedicated resource name.

* `engine_name` - The API name.

* `availability_zone` - The AZ information.

* `architecture` - The type of the compute host in the dedicated resource.

* `capacity` - The capacity of the dedicated resource.
  The [capacity](#capacity_struct) structure is documented below.

* `status` - The status of a dedicated resource.

<a name="capacity_struct"></a>
The `capacity` block supports:

* `vcpus` - The CPU cores.

* `ram` - The memory size, in GB.

* `volume` - The storage size, in GB.
