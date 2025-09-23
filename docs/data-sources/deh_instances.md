---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_instances"
description: |-
  Use this data source to get the list of dedicated host instances.
---

# huaweicloud_deh_instances

Use this data source to get the list of dedicated host instances.

## Example Usage

```hcl
data "huaweicloud_deh_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `dedicated_host_id` - (Optional, String) Specifies the ID of the dedicated host.

* `name` - (Optional, String) Specifies the name of the dedicated host.

* `host_type` - (Optional, String) Specifies the type of the dedicated host.

* `host_type_name` - (Optional, String) Specifies the name of the dedicated host type.

* `flavor` - (Optional, String) Specifies the flavor ID.

* `state` - (Optional, String) Specifies the status of the dedicated host.
  Value options: **available**, **fault** or **released**.

* `availability_zone` - (Optional, String) Specifies the AZ to which the dedicated host belongs.

* `tags` - (Optional, String) Specifies the tags of the dedicated host.

* `instance_uuid` - (Optional, String) Specifies the ID of the ECS on the dedicated host.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dedicated_hosts` - Indicates the list of dedicated hosts.

  The [dedicated_hosts](#dedicated_hosts_struct) structure is documented below.

<a name="dedicated_hosts_struct"></a>
The `dedicated_hosts` block supports:

* `dedicated_host_id` - Indicates the ID of the dedicated host.

* `name` - Indicates the name of the dedicated host.

* `auto_placement` - Indicates whether to allow an ECS to be placed on any available dedicated host
  if its dedicated host ID is not specified during its creation.

* `availability_zone` - Indicates the AZ to which the dedicated host belongs.

* `state` - Indicates the status of the dedicated host.

* `available_memory` - Indicates the available memory size of the dedicated host.

* `instance_total` - Indicates the total number of ECSs on the dedicated host.

* `instance_uuids` - Indicates the UUIDs of the ECSs running on the dedicated host.

* `available_vcpus` - Indicates the number of available vCPUs for the dedicated host.

* `host_properties` - Indicates the properties of the dedicated host.

  The [host_properties](#dedicated_hosts_host_properties_struct) structure is documented below.

* `tags` - Indicates the tags of the dedicated host.

* `sys_tags` - Indicates the system tags of the dedicated host.

* `allocated_at` - Indicates the time when the dedicated host is allocated.

<a name="dedicated_hosts_host_properties_struct"></a>
The `host_properties` block supports:

* `host_type` - Indicates the type of the dedicated host.

* `host_type_name` - Indicates the name of the dedicated host type.

* `vcpus` - Indicates the number of vCPUs on the dedicated host.

* `cores` - Indicates the number of physical cores on the dedicated host.

* `sockets` - Indicates the number of physical sockets on the dedicated host.

* `memory` - Indicates the size of physical memory on the dedicated host.

* `available_instance_capacities` - Indicates the flavors of ECSs placed on the dedicated host.

  The [available_instance_capacities](#host_properties_available_instance_capacities_struct) structure is documented below.

<a name="host_properties_available_instance_capacities_struct"></a>
The `available_instance_capacities` block supports:

* `flavor` - Indicates the specifications of ECSs that can be created.
