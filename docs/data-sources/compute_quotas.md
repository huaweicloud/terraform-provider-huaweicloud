---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_quotas"
description: |-
  Use this data source to get the quotas of all resources for a specified tenant, including used quotas.
---

# huaweicloud_compute_quotas

Use this data source to get the quotas of all resources for a specified tenant, including used quotas.

## Example Usage

```hcl
data "huaweicloud_compute_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `absolute` - Indicates the tenant quotas.

  The [absolute](#absolute_struct) structure is documented below.

<a name="absolute_struct"></a>
The `absolute` block supports:

* `max_security_groups` - Indicates the maximum number of security groups you can use.
  The quota complies with the VPC quota limit.

* `max_server_group_members` - Indicates the maximum number of ECSs in an ECS group.

* `max_server_groups` - Indicates the maximum number of server groups.

* `max_total_cores` - Indicates the maximum number of CPU cores that the current tenant can apply for.

* `max_total_floating_ips` - Indicates the maximum number of floating IP addresses you can use.

* `max_total_spot_instances` - Indicates the maximum number of Spot Instances that can be applied for.

* `max_total_spot_ram_size` - Indicates the maximum memory size (MiB) allowed.

* `total_spot_cores_used` - Indicates the number of CPU cores currently used by the spot instance.

* `max_personality` - Indicates the maximum number of files that can be injected.

* `max_security_group_rules` - Indicates the maximum number of security group rules that you can configure in a security
  group. The quota complies with the VPC quota limit.

* `total_security_groups_used` - Indicates the number of used security groups.

* `total_spot_ram_used` - Indicates the memory usage of the current spot instance (unit: MB).

* `max_server_meta` - Indicates the maximum length of the metadata you can use.

* `max_total_instances` - Indicates the maximum number of ECSs that can be requested.

* `max_total_ram_size` - Indicates the maximum memory size (MiB) allowed.

* `total_cores_used` - Indicates the number of the used CPU cores.

* `total_instances_used` - Indicates the number of used ECSs.

* `total_ram_used` - Indicates the used memory size (MiB).

* `max_image_meta` - Indicates the maximum length of the image metadata.

* `max_personality_size` - Indicates the maximum size (byte) of the file to be injected.

* `max_total_keypairs` - Indicates the maximum number of SSH key pairs you can use.

* `total_floating_ips_used` - Indicates the number of used floating IP addresses.

* `total_server_groups_used` - Indicates the number of used server groups.

* `max_total_spot_cores` - Indicates the maximum number of CPU cores that can be requested for a spot instance.

* `total_spot_instances_used` - Indicates the current number of Spot instances in use.

* `max_cluster_server_group_members` - Indicates the maximum number of cluster server group.

* `max_fault_domain_members` - Indicates the maximum number of fault domain.
