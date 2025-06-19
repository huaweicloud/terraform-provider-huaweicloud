---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_quotas"
description: |-
  Use this data source to query SDRS quotas within HuaweiCloud.
---

# huaweicloud_sdrs_quotas

Use this data source to query SDRS quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The tenant quota information.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The resource quota list.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `max` - The maximally allowed resource quota. Value `-1` means to unlimited.

* `type` - The resource type. Valid values are:
  + **server_groups**: Indicates protection groups.
  + **replications**: Indicates replication pairs.

* `used` - The number of used resources.

* `quota` - The resource quota. Value `-1` means to unlimited.

* `min` - The minimally allowed resource quota.
