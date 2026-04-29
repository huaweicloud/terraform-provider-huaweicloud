---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_public_ip_pool_types"
description: |-
  Use this dataSource to get the list of public IP pool types.
---

# huaweicloud_vpc_eip_public_ip_pool_types

Use this dataSource to get the list of public IP pool types.

## Example Usage

```hcl
data "huaweicloud_vpc_eip_public_ip_pool_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fields` - (Optional, List) Specifies the fields returned by the query.  
  The supported field names include: **id**, **name**, **size**, **used**, **project_id**, **status**,
  **billing_info**, **created_at**, **updated_at**, **type**, **shared**, **is_common**,
  **description**, **tags**, **enterprise_project_id**, **allow_share_bandwidth_types**, and **public_border_group**.

* `sort_key` - (Optional, String) Specifies the sort key. The supported field names are **id**, **name**,
  **created_at**, **updated_at**, and **public_border_group**.

* `sort_dir` - (Optional, String) Specifies the sort order. Valid values are **asc** and **desc**.

* `type_id` - (Optional, String) Specifies the type ID.

* `name` - (Optional, String) Specifies the name

* `size` - (Optional, Int) Specifies the size

* `status` - (Optional, String) Specifies the status.

* `type` - (Optional, String) Specifies the type.

* `description` - (Optional, String) Specifies the description.

* `public_border_group` - (Optional, String) Specifies the public_border_group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `public_ip_pool_types` - The list of public IP pool type objects.

  The [public_ip_pool_types](#public_ip_pool_types_struct) structure is documented below.

<a name="public_ip_pool_types_struct"></a>
The `public_ip_pool_types` block supports:

* `id` - The public IP pool type ID.

* `type` - The public IP pool type.
