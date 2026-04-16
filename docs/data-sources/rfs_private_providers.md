---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_providers"
description: |-
  Use this datasource to get the list of private providers.
---

# huaweicloud_rfs_private_providers

Use this datasource to get the list of private providers.

## Example Usage

```hcl
data "huaweicloud_rfs_private_providers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `sort_key` - (Optional, String) Specifies the sort field, only supports **create_time**.

* `sort_dir` - (Optional, String) Specifies the ascending or descending order.
  Valid values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `providers` - The list of private providers.

  The [providers](#providers_struct) structure is documented below.

<a name="providers_struct"></a>
The `providers` block supports:

* `provider_id` - The unique ID of the private provider.

* `provider_name` - The name of the private provider.

* `provider_description` - The description of the private provider.

* `provider_source` - The source parameter that users need to specify when defining the required providers information
  in Terraform templates using private providers.

* `provider_agency_urn` - The IAM agency URN bound to the custom provider.

* `provider_agency_name` - The IAM agency name bound to the custom provider.

* `create_time` - The creation time of a private provider. It is represented in UTC format (YYYY-MM-DDTHH:mm:ss.SSSZ),
  such as **1970-01-01T00:00:00.000Z**.

* `update_time` - The update time of a private provider. It is represented in UTC format (YYYY-MM-DDTHH:mm:ss.SSSZ),
  such as **1970-01-01T00:00:00.000Z**.
