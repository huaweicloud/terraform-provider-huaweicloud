---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_segments"
description: |-
  Use this data source to get the list of global EIP segments.
---

# huaweicloud_global_eip_segments

Use this data source to get the list of global EIP segments.

## Example Usage

```hcl
data "huaweicloud_global_eip_segments" "test" {}
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
  Supported values include **id**, **name**, **description**, **domain_id**, **access_site**, **geip_pool_name**,
  **isp**, **ip_version**, **cidr**, **cidr_v6**, **freezen**, **freezen_info**, **status**, **created_at**,
  **updated_at**, **internet_bandwidth**, **associate_instance**, **is_pre_paid**, and **enterprise_project_id**.

* `sort_key` - (Optional, String) Specifies the sort field.

* `sort_dir` - (Optional, String) Specifies the sort directions.  
  Valid values are **asc** and **desc**.

* `segment_ids` - (Optional, List) Specifies the global EIP segment IDs to filter.

* `internet_bandwidth_id` - (Optional, List) Specifies the global internet bandwidth IDs to filter.

* `name` - (Optional, List) Specifies the segment names to filter.

* `name_like` - (Optional, String) Specifies the fuzzy name match string to filter.

* `access_site` - (Optional, List) Specifies the access sites to filter.

* `geip_pool_name` - (Optional, List) Specifies the global EIP pool names to filter.

* `isp` - (Optional, List) Specifies the ISP lines to filter.

* `ip_address` - (Optional, List) Specifies the IPv4 addresses used to match CIDR to filter.

* `ipv6_address` - (Optional, List) Specifies the IPv6 addresses used to match CIDR_V6 to filter.

* `ip_version` - (Optional, List) Specifies the IP versions to filter.

* `cidr` - (Optional, List) Specifies the IPv4 CIDR blocks to filter.

* `cidr_v6` - (Optional, List) Specifies the IPv6 CIDR blocks to filter.

* `freezen` - (Optional, List) Specifies whether the segments are frozen to filter.

* `internet_bandwidth_is_null` - (Optional, List) Specifies whether the segments are bound to a global internet
  bandwidth to filter.

* `status` - (Optional, List) Specifies the segment statuses to filter.  
  The valid values are as follows:
  + **IDLE**: The global elastic public IP is not bound to an instance.
  + **INUSE**: Global elastic public IP in use.
  + **FREEZED**: The global elastic public IP has been frozen.
  + **PENDING_CREATE**: Global elastic public IP is being created.
  + **PENDING_UPDATE**: The global elastic public IP is being updated.

* `associate_instance_region` - (Optional, List) Specifies the regions of bound instances to filter.

* `associate_instance_instance_type` - (Optional, List) Specifies the bound instance types to filter.
  Valid values include **DC-CONNECT-GATEWAY** and **IPV6-DC-CONNECT-GATEWAY**.

* `associate_instance_public_border_group` - (Optional, List) Specifies the public border groups of bound instances to
  filter.

* `associate_instance_instance_site` - (Optional, List) Specifies the sites of bound instances to filter.

* `associate_instance_instance_id` - (Optional, List) Specifies the bound instance IDs to filter.

* `associate_instance_project_id` - (Optional, List) Specifies the project IDs of bound instances to filter.

* `associate_instance_service_id` - (Optional, List) Specifies the service IDs of bound instances to filter.

* `associate_instance_service_type` - (Optional, List) Specifies the service types of bound instances to filter.

* `enterprise_project_id` - (Optional, List) Specifies the enterprise project IDs to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `global_eip_segments` - The list of global EIP segments.

  The [global_eip_segments](#global_eip_segments_struct) structure is documented below.

<a name="global_eip_segments_struct"></a>
The `global_eip_segments` block supports:

* `id` - The global EIP segment ID.

* `name` - The segment name.

* `description` - The segment description.

* `domain_id` - The domain ID.

* `access_site` - The access site.

* `geip_pool_name` - The global EIP pool name.

* `isp` - The ISP line.

* `ip_version` - The IP version.

* `cidr` - The IPv4 CIDR block.

* `cidr_v6` - The IPv6 CIDR block.

* `freezen` - Whether the segment is frozen.

* `freezen_info` - The freeze information.

* `status` - The segment status.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `internet_bandwidth` - The internet bandwidth bound to the segment.

  The [internet_bandwidth](#internet_bandwidth_struct) structure is documented below.

* `is_pre_paid` - Whether the segment is prepaid.

* `is_charged` - Whether the segment is charged.

* `enterprise_project_id` - The enterprise project ID.

<a name="internet_bandwidth_struct"></a>
The `internet_bandwidth` block supports:

* `id` - The internet bandwidth ID.

* `size` - The bandwidth size.
