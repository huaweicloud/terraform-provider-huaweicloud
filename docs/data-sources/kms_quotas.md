---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_quotas"
description: |-
  Use this data source to get the quotas of KMS within HuaweiCloud.
---

# huaweicloud_kms_quotas

Use this data source to get the quotas of KMS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_kms_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota details.
  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The list of the resource quotas.
  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `type` - The type of the quotas. The valid values are as follows:
  + **CMK**: The user master key.
  + **grant_per_CMK**: The number of authorizations a user master key can create.

* `used` - The number of quotas used.

* `quota` - The total number of quotas.
