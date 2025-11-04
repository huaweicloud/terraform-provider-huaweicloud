---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_retirable_grants"
description: |-
  Use this datasource to get the list of retirable grants.
---

# huaweicloud_kms_retirable_grants

Use this datasource to get the list of retirable grants.

## Example Usage

```hcl
data "huaweicloud_kms_retirable_grants" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `sequence` - (Optional, String) Specifies the request sequence number.
  For example, **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `grants` - The list of the retirable grants.

  The [grants](#grants_struct) structure is documented below.

<a name="grants_struct"></a>
The `grants` block supports:

* `key_id` - The key ID.

* `grant_id` - The grant ID.

* `name` - The grant name.

* `grantee_principal` - The ID of the authorized user or account.

* `grantee_principal_type` - The authorization type.
  The value can be **user** or **domain**.

* `operations` - The list of granted operations.
  The valid values are as follows:
  + **create-datakey**
  + **create-datakey-without-plaintext**
  + **encrypt-datakey**
  + **decrypt-datakey**
  + **describe-key**
  + **retire-grant**
  + **encrypt-data**
  + **decrypt-data**

* `issuing_principal` - The ID of the user who created the grant.

* `retiring_principal` - The ID of the user who retirable the grant.

* `creation_date` - The creation time.
