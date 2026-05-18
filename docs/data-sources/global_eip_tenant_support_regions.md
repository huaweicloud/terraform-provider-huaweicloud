---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_tenant_support_regions"
description: |-
  Use this data source to get the list of regions that can be bound for global EIP at the given access site.
---

# huaweicloud_global_eip_tenant_support_regions

Use this data source to get the list of regions that can be bound for global EIP at the given access site.

## Example Usage

```hcl
variable "access_site" {}

data "huaweicloud_global_eip_tenant_support_regions" "test" {
  access_site = var.access_site
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `access_site` - (Required, String) Specifies the access site.  
  The value can be obtained using the `huaweicloud_global_eip_access_sites` data source.

* `fields` - (Optional, List) Specifies the fields to return.  
  Supported values include **id**, **instance_type**, **region_id**, **public_border_group**, **access_site**,
  **status**, **created_at**, and **updated_at**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `support_regions` - The list of global EIP support region objects for the site.

  The [support_regions](#support_regions_struct) structure is documented below.

<a name="support_regions_struct"></a>
The `support_regions` block supports:

* `id` - The ID of the region binding record.

* `instance_type` - The supported instance type.  
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

* `public_border_group` - Indicates whether the resource is in the central site or at an edge site.
  The value is **center** or an edge site name.

* `region_id` - The region ID.

* `access_site` - The access site.

* `status` - The site status.  
  The valid values are as follows:
  + **ACTIVE**: Already online.
  + **INACTIVE**: Offline.

* `created_at` - The creation time.

* `updated_at` - The update time.
