---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_instance"
description: |-
  Manages a DEH instance resource within HuaweiCloud.
---

# huaweicloud_deh_instance

Manages a DEH instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "availability_zone" {}
variable "name" {}
variable "host_type" {}

resource "huaweicloud_deh_instance" "test" {
  availability_zone = var.availability_zone
  name              = "deh_name"
  host_type         = var.host_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the AZ to which the dedicated host belongs.

* `name` - (Required, String) Specifies the name of the dedicated host. It can contain a maximum of 255 characters and
  cannot start or end with spaces.

* `host_type` - (Required, String, NonUpdatable) Specifies the type of the dedicated host.

* `auto_placement` - (Optional, String) Specifies whether to allow an ECS to be placed on any available dedicated host if
  its dedicated host ID is not specified during its creation. Value options: **on** and **off**.

* `metadata` - (Optional, Map, NonUpdatable) Specifies the metadata of the dedicated host.

* `charging_mode` - (Optional, String, NonUpdatable) Specifies the charging mode of the dedicated host. Value options:
  **prePaid**.

* `period_unit` - (Optional, String, NonUpdatable) Specifies the charging period unit. Value options: **month** and
  **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, String, NonUpdatable) Specifies the charging period.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Value options: **true** and **false**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the dedicated host.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `host_properties` - Indicates the properties of the dedicated host.
  The [host_properties](#host_properties_struct) structure is documented below.

* `state` - Indicates the status of the dedicated host.

* `available_vcpus` - Indicates the number of available vCPUs of dedicated host.

* `available_memory` - Indicates the available memory size of dedicated host.

* `allocated_at` - Indicates the time when the dedicated host is allocated.

* `instance_total` - Indicates the total number of ECSs on the dedicated host.

* `instance_uuids` - Indicates the UUIDs of the ECSs running on the dedicated host.

* `sys_tags` - Indicates the system tags of the dedicated host.

<a name="host_properties_struct"></a>
The `host_properties` block supports:

* `host_type` - Indicates the type of the dedicated host.

* `host_type_name` - Indicates the name of the dedicated host type.

* `vcpus` - Indicates the number of vCPUs on the dedicated host.

* `cores` - Indicates the number of physical cores on the dedicated host.

* `sockets` - Indicates the number of physical sockets on the dedicated host.

* `memory` - Indicates the size of physical memory on the dedicated host.

* `available_instance_capacities` - Indicates the flavors of ECSs placed on the dedicated host.
  The [available_instance_capacities](#available_instance_capacities_struct) structure is documented below.

<a name="available_instance_capacities_struct"></a>
The `available_instance_capacities` block supports:

* `flavor` - Indicates the specifications of ECSs that can be created.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The DEH instance can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_deh_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `period_unit`,`period`, `auto_renew` and
`enterprise_project_id`. It is generally recommended running `terraform plan` after importing a DEH instance. You can
then decide if changes should be applied to the DEH instance, or the resource definition should be updated to align with
the mesh. Also you can ignore changes as below.

```hcl
resource "huaweicloud_deh_instance" "test" {
    ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew, enterprise_project_id,
    ]
  }
}
```
