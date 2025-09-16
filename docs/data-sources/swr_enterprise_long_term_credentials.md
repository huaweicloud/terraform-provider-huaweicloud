---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_long_term_credentials"
description: |-
  Use this data source to get the list of SWR enterprise instance long term credentials.
---

# huaweicloud_swr_enterprise_long_term_credentials

Use this data source to get the list of SWR enterprise instance long term credentials.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_long_term_credentials" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `auth_tokens` - Indicates the namespaces.
  The [auth_tokens](#attrblock--auth_tokens) structure is documented below.

<a name="attrblock--auth_tokens"></a>
The `auth_tokens` block supports:

* `id` - Indicates the namespace ID.

* `name` - Indicates the credential name.

* `enable` - Indicates whether to enable the credential.

* `user_id` - Indicates the user ID.

* `user_profile` - Indicates the user profile.

* `created_at` - Indicates the creation time.

* `expire_date` - Indicates the expired time.
