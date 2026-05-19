---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_quotas"
description: |-
  Use this data source to get the quotas information of DRS within HuaweiCloud.
---

# huaweicloud_drs_quotas

Use this data source to get the quotas information of DRS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_quotas" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quotas information.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resource` - The resource quotas information.

  The [resource](#resource_struct) structure is documented below.

<a name="resource_struct"></a>
The `resource` block supports:

* `type` - The quota type.
  The valid values are as follows:
  + **instances**: The number of instances.
  + **cpu**: The CPU quota.
  + **cores**: The CPU cores quota.
  + **server_groups**: The server groups quota.
  + **mem**: The memory quota.
  + **ram**: The RAM quota.

* `min` - The minimum value of the quota.

* `max` - The maximum value of the quota.

* `quota` - The actual value of the user quota.

* `used` - The used value of the quota.
