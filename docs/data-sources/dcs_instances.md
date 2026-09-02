---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instances"
description: |-
  Use this data source to get the list of DCS instances.
---

# huaweicloud_dcs_instances

Use this data source to get the list of DCS instances.

## Example Usage

### Query all instances

```hcl
data "huaweicloud_dcs_instances" "test" {}
```

### Filter by name

```hcl
variable "instance_name" {}

data "huaweicloud_dcs_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of an instance.

* `status` - (Optional, String) Specifies the cache instance status.
  The valid values are **RUNNING**, **ERROR**, **RESTARTING**, **FROZEN**, **EXTENDING**, **RESTORING**, **FLUSHING**.

* `private_ip` - (Optional, String) Specifies the subnet Network ID.

* `capacity` - (Optional, Float) Specifies the cache capacity. Unit: GB.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `tags` - (Optional, String) Specifies the tags of the instance.
  If multiple tag key-value pairs are used for the query at the same time, they should be separated by commas(,),
  indicating that the query includes instances that have the specified tag key-value pairs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of DCS instances.
  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

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

* `publicip_id` - Indicates the EIP ID associated with the instance.

* `created_at` - Indicates the creation time of the instance.

* `updated_at` - Indicates the update time of the instance.

* `enable_ssl` - Whether SSL is enabled for the instance.

* `publicip_address` - Indicates the public IP address associate with the instance.

* `service_upgrade` - Whether an upgrade task has been created for the instance.

* `no_password_access` - Whether password-protected access is enabled for the instance.

* `service_task_id` - Indicates the ID of the upgrade task.

* `user_id` - Indicates the ID of the user to which the instance belongs.

* `user_name` - Indicates the username of the instance.

* `readonly_domain_name` - Indicates the read-only domain name of the instance.

* `cpu_type` - Indicates the CPU type of the instance.

* `az_codes` - Indicates the availability zones where there are available resources.
