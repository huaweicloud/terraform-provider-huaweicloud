---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instances"
description: ""
---

# huaweicloud_dcs_instances

Use this data source to get the list of DCS instances.

## Example Usage

```hcl
data "huaweicloud_dcs_instances" "test" {
  name   = "test_name"
  status = "RUNNING"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of an instance.

* `status` - (Optional, String) Specifies the cache instance status. The valid values are **RUNNING**, **ERROR**,
  **RESTARTING**, **FROZEN**, **EXTENDING**, **RESTORING**, **FLUSHING**.

* `private_ip` - (Optional, String) Specifies the subnet Network ID.

* `capacity` - (Optional, Float) Specifies the cache capacity. Unit: GB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of DCS instances.
  The [Instance](#DcsInstance_Instance) structure is documented below.

<a name="DcsInstance_Instance"></a>
The `Instance` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the name of an instance.

* `engine` - Indicates a cache engine.

* `engine_version` - Indicates the version of a cache engine.

* `capacity` - Indicates the cache capacity. Unit: GB.

* `flavor` - Indicates the flavor of the cache instance.

* `availability_zones` - Specifies the code of the AZ where the cache node resides.

* `vpc_id` - Indicates the ID of VPC which the instance belongs to.

* `vpc_name` - Indicates the name of VPC which the instance belongs to.

* `subnet_id` - Indicates the ID of subnet which the instance belongs to.

* `subnet_name` - Indicates the name of subnet which the instance belongs to.

* `security_group_id` - Indicates the ID of the security group which the instance belongs to.

* `security_group_name` - Indicates the name of security group which the instance belongs to.

* `enterprise_project_id` - Indicates the enterprise project id of the dcs instance.

* `description` - Indicates the description of an instance.

* `private_ip` - Indicates the IP address of the DCS instance.

* `maintain_begin` - Indicates the time at which the maintenance time window starts.

* `maintain_end` - Indicates the time at which the maintenance time window ends.

* `charging_mode` - Indicates the charging mode of the cache instance.

* `port` - Indicates the port of the cache instance.

* `status` - Indicates the cache instance status.

* `used_memory` - Indicates the size of the used memory. Unit: MB.

* `max_memory` - Indicates the total memory size. Unit: MB.

* `domain_name` - Indicates the domain name of the instance.

* `access_user` - Indicates the username used for accessing a DCS Memcached instance.

* `order_id` - Indicates the ID of the order that created the instance.

* `tags` - Indicates The key/value pairs to associate with the DCS instance.
