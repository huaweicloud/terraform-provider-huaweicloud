---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_registries"
description: |-
  Use this data source to get the list of SWR enterprise instance registries.
---

# huaweicloud_swr_enterprise_instance_registries

Use this data source to get the list of SWR enterprise instance registries.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_instance_registries" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `name` - (Optional, String) Specifies the name.

* `type` - (Optional, String) Specifies the repository type

* `order_column` - (Optional, String) Specifies the order column.
  Values can be **created_at** or **updated_at**. Default to **created_at**.

* `order_type` - (Optional, String) Specifies the order type. Values can be **desc** or **asc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `registries` - Indicates the repositories.

  The [registries](#registries_struct) structure is documented below.

* `total` - Indicates the total number of the repositories.

<a name="registries_struct"></a>
The `registries` block supports:

* `id` - Indicates the repository ID.

* `name` - Indicates the repository name.

* `type` - Indicates the repository type.

* `url` - Indicates the repository url.

* `insecure` - Indicates whether to verify the remote certificate.

* `description` - Indicates the repository description.

* `region_id` - Indicates the region ID of the target instance.

* `instance_id` - Indicates the ID of the target instance.

* `credential` - Indicates the credential infos.

  The [credential](#registries_credential_struct) structure is documented below.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.

* `status` - Indicates the repository status.

<a name="registries_credential_struct"></a>
The `credential` block supports:

* `access_key` - Indicates the access key.

* `access_secret` - Indicates the access secret.

* `type` - Indicates the credential type. Value can be **basic**.
