---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_cloud_connection_quotas"
description: |-
  Use this data source to get the cloud connection resource quotas.
---

# huaweicloud_cc_cloud_connection_quotas

Use this data source to get the cloud connection resource quotas.

## Example Usage

```hcl
data "huaweicloud_cc_cloud_connection_quotas" "test"{
  quota_type = "cloud_connection"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `quota_type` - (Required, String) Specifies the quota type.
  Value options:
  + **cloud_connection**: the maximum number of cloud connections that can be created in an account
  + **cloud_connection_region**: the maximum number of regions where a cloud connection can be used
  + **cloud_connection_route**: the maximum number of routes that can be added to a cloud connection
  + **region_network_instance**: the maximum number of network instances that can be loaded to a cloud connection in a region

* `cloud_connection_id` - (Optional, String) Specifies the cloud connection ID.
  This parameter is mandatory when you query the value of each of the three parameters:
  **cloud_connection_region**, **cloud_connection_route**, **and region_network_instance**.

* `region_id` - (Optional, String) Specifies the region ID.
  This parameter is mandatory when you query the value of **region_network_instance**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the quota list.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `cloud_connection_id` - Indicates the cloud connection ID.

* `region_id` - Indicates the region ID.

* `quota_type` - Indicates the quota type.

* `quota_number` - Indicates the total quotas.

* `quota_used` - Indicates the used quotas.
