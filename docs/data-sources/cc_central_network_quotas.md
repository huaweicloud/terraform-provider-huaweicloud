---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_quotas"
description: |-
  Use this data source to get the list of CC central network quotas.
---

# huaweicloud_cc_central_network_quotas

Use this data source to get the list of CC central network quotas.

## Example Usage

```hcl
data "huaweicloud_cc_central_network_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `quota_type` - (Optional, List) Specifies the quota type.
  The valid values are as follows:
  + **central_networks_per_account**
  + **policy_versions_per_central_network**
  + **size_of_document_per_central_network_policy_version**
  + **planes_per_central_network**
  + **er_instances_per_region_per_central_network**
  + **connections_per_central_network**
  + **attachments_per_central_network**
  + **GDGW_attachments_per_region_per_central_network**
  + **ER_ROUTE_TABLE_attachments_per_region_per_central_network**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota list.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `quota_key` - The central network quota type.

* `quota_limit` - The quota size.

* `used` - The used quotas.

* `unit` - The unit of the quota value.
