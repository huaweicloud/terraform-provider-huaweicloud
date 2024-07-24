---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_grants"
description: |-
  Use this datasource to get a list of grants.
---

# huaweicloud_kms_grants

Use this datasource to get a list of grants.

## Example Usage

```hcl
variable "key_id" {}

data "huaweicloud_kms_grants" "test" {
  key_id = var.key_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `key_id` - (Required, String) Specifies the key ID to which the grants belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `grants` - The list of the grants.

  The [grants](#grants_struct) structure is documented below.

<a name="grants_struct"></a>
The `grants` block supports:

* `id` - The ID of the grant.

* `key_id` - The key ID to which the grant belongs.

* `name` - The name of the grant.

* `type` - The authorization type.

* `creator` - The ID of the user who created the grant.

* `grantee_principal` - The ID of the authorized user or account.

* `operations` - List of granted operations.

* `created_at` - The creation time of the grant, in RFC3339 format.
