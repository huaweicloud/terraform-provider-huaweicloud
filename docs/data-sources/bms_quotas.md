---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_quotas"
description: |-
  Use this data source to get the quotas of all resources for a specified tenant, including used quotas.
---

# huaweicloud_bms_quotas

Use this data source to get the quotas of all resources for a specified tenant, including used quotas.

## Example Usage

```hcl
data "huaweicloud_bms_quotas" "demo" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `absolute` - Indicates tenant quotas.

  The [absolute](#absolute_struct) structure is documented below.

<a name="absolute_struct"></a>
The `absolute` block supports:

* `max_total_instances` - Indicates the maximum number of BMSs you can create.

* `max_server_group_members` - Indicates the maximum number of BMSs in a server group.

* `total_server_groups_used` - Indicates the number of used server groups.

* `max_security_groups` - Indicates the maximum number of security groups you can use.

* `max_image_meta` - Indicates the maximum length of the image metadata.

* `total_cores_used` - Indicates the number of used CPU cores.

* `max_total_keypairs` - Indicates the maximum number of SSH key pairs you can use.

* `max_personality` - Indicates the maximum number of files that can be injected.

* `max_security_group_rules` - Indicates the maximum number of security group rules that you can configure in a security
  group.

* `max_total_floating_ips` - Indicates the maximum number of EIPs you can use.

* `total_instances_used` - Indicates the number of the used BMSs.

* `max_total_ram_size` - Indicates the maximum memory (MB) you can use.

* `max_server_meta` - Indicates the maximum length of the metadata you can use.

* `max_personality_size` - Indicates the maximum size (byte) of the file to be injected.

* `max_server_groups` - Indicates the maximum number of server groups.

* `total_floating_ips_used` - Indicates the number of used EIPs.

* `max_total_cores` - Indicates the maximum number of CPU cores you can use.

* `total_ram_used` - Indicates the used memory (MB).

* `total_security_groups_used` - Indicates the number of used security groups.
