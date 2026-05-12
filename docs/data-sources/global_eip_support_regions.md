---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_support_regions"
description: |-
  Use this data source to get the list of regions that global EIPs can bind to.
---

# huaweicloud_global_eip_support_regions

Use this data source to get the list of regions that global EIPs can bind to.

## Example Usage

```hcl
data "huaweicloud_global_eip_support_regions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `page_reverse` - (Optional, String) Specifies the page direction.  
  The valid values are as follows:
  + **true**: means the previous page.
  + **false**: means the next page.

* `fields` - (Optional, List) Specifies the fields to return.
  Supported values include **id**, **instance_type**, **access_site**, **region_id**, **public_border_group**,
  **remote_endpoint**, **status**, **created_at**, and **updated_at**.

* `sort_key` - (Optional, String) Specifies the sort fields.

* `sort_dir` - (Optional, String) Specifies the sort directions.
  Valid values are **asc** and **desc**.

* `support_region_ids` - (Optional, List) Specifies the support region record IDs to filter.

* `instance_type` - (Optional, List) Specifies the supported instance types to filter.  
  The valid values are as follows:
  + **DC-CONNECT-GATEWAY**
  + **IPV6-DC-CONNECT-GATEWAY**
  + **ECS**
  + **IPV6-ECS**
  + **PORT**
  + **IPV6-PORT**
  + **VIP**
  + **IPV6-VIP**
  + **ELB**
  + **IPV6-ELB**
  + **NATGW**

* `public_border_group` - (Optional, List) Specifies the public border groups to filter. Valid value is **center** or an
  edge site name.

* `access_site` - (Optional, List) Specifies the access sites to filter.

* `region_id` - (Optional, List) Specifies the region IDs to filter.

* `remote_endpoint` - (Optional, List) Specifies the remote endpoints to filter.

* `status` - (Optional, List) Specifies the values used to filter by EIP or site status in the query.  
  The valid values are as follows:
  + **IDLE**: The global elastic public IP is not bound to an instance.
  + **INUSE**: Global elastic public IP in use.
  + **FREEZED**: The global elastic public IP has been frozen.
  + **PENDING_CREATE**: Global elastic public IP is being created.
  + **PENDING_UPDATE**: The global elastic public IP is being updated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `support_regions` - The list of global EIP support region objects.

  The [support_regions](#support_regions_struct) structure is documented below.

<a name="support_regions_struct"></a>
The `support_regions` block supports:

* `id` - The support region record ID.

* `instance_type` - The supported instance type.

* `access_site` - The access site.

* `region_id` - The region ID.

* `public_border_group` - The public border group.

* `remote_endpoint` - The next-hop address.

* `status` - The status.  
  The valid values are as follows:
  + **ACTIVE**: Already online.
  + **INACTIVE**: Offline.

* `created_at` - The creation time.

* `updated_at` - The update time.
