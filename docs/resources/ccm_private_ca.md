---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca"
description: |
  Manages CCM private CA resource within HuaweiCloud.
---

# huaweicloud_ccm_private_ca

Manages CCM private CA resource within HuaweiCloud.

## Example Usage

### create a root private CA

```hcl
variable "obs_bucket_name" {}

resource "huaweicloud_ccm_private_ca" "test_root" {
  region              = "cn-north-4"
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"

  distinguished_name {
    common_name         = "test-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 5
  }

  crl_configuration {
    obs_bucket_name = var.obs_bucket_name
    valid_days      = "7"
    enabled         = true
  }
}
```

### create a subordinate private CA

```hcl
variable "root_issuer_id" {}

resource "huaweicloud_ccm_private_ca" "test_subordinate" {
  region              = "cn-north-4"
  type                = "SUBORDINATE"
  key_algorithm       = "RSA2048"
  issuer_id           = var.root_issuer_id
  signature_algorithm = "SHA512"
  pending_days        = "7"

  distinguished_name {
    common_name         = "test-subordinate"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  
  validity {
    type  = "DAY"
    value = 4
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCM private CA. If omitted, the
  provider-level region will be used. Changing this will create a new resource. Now only support cn-north-4 (china) and
  ap-southeast-3 (international).

* `type` - (Required, String, ForceNew) Specifies the type of private CA. Options are: **ROOT**, **SUBORDINATE**.
  Changing this parameter will create a new resource.

* `distinguished_name` - (Required, List, ForceNew) Specifies the distinguished name of private CA.
  Changing this parameter will create a new resource.
  The [distinguished_name](#block-distinguished_name) structure is documented below.

* `key_algorithm` - (Required, String, ForceNew) Specifies the key algorithm of private CA.
  Options are: **RSA2048**, **RSA4096**, **EC256**, **EC256**,**SM2**.
  Changing this parameter will create a new resource.

* `signature_algorithm` - (Required, String, ForceNew) Specifies the signature algorithm of private CA.
  Options are: **SHA256**, **SHA384**, **SHA512**, **SM3**.
  Changing this parameter will create a new resource.

* `validity` - (Required, List, ForceNew) Specifies the validity of private CA.
  Changing this parameter will create a new resource.
  The [validity](#block-validity) structure is documented below.

* `pending_days` - (Required, String, ForceNew) Specifies the pending days when deleting the private CA. It's limited
  between `7` to `30`. Changing this parameter will create a new resource.

* `issuer_id` - (Optional, String, ForceNew) Specifies the ID of the parent CA. It's **required** for subordinate CA.
  Changing this parameter will create a new resource.

* `path_length` - (Optional, Int, ForceNew) Specifies the length of the CA certificate path. The valid value is
  limited between `0` to `6`. If you want to create a root CA, this parameter is **not required** by default and the
  value will be set to `7` in return. Changing this parameter will create a new resource.

* `key_usages` - (Optional, List, ForceNew) Specifies the key usage of private CA. It's a list of string.
  Options are: **digitalSignature**, **nonRepudiation**, **keyEncipherment**, **dataEncipherment**, **keyAgreement**,
  **keyCertSign**, **cRLSign**, **encipherOnly**, **decipherOnly**.
  This parameter is [digitalSignature,keyCertSign,cRLSign] by default and only support to customize when you create a
  subordinate CA. Changing this parameter will create a new resource.

* `crl_configuration` - (Optional, List) Specifies the CRL configuration of private CA. The maximum length is `1`.
  The [crl_configuration](#block-crl_configuration) structure is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this parameter will create a new resource.

* `action` - (Optional, String) Specifies the operation action of the CCM private CA.
  Valid values are **enable** and **disable**.

  -> 1. The **enable** operation is allowed only when the CA status is **DISABLED**.
  <br/>2. The **disable** operation is allowed only when the CA status is **ACTIVED** or **EXPIRED**.
  <br/>3. The status of the newly created CA is **ACTIVED**.
  <br/>4. The child CA's certificate status will remain consistent with the parent CA, so enabling or disabling the child
  CA alone may not be effective.

* `charging_mode` - (Optional, String, ForceNew) Specifies the billing mode of the private CA.
  The valid values are **prePaid** and **postPaid**. Defaults to **postPaid**.
  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto-renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs associating with the private CA.

<a name="block-distinguished_name"></a>
The `distinguished_name` block supports:

* `common_name` - (Required, String, ForceNew) Specifies the common name of private CA. The valid length is limited
  between `1` to `64`, Only Chinese and English letters, digits, hyphens (-), underscores (_), dots (.), comma (,),
  space ( ) and asterisks (*) are allowed. Changing this parameter will create a new resource.

* `country` - (Required, String, ForceNew) Specifies the country of private CA. The valid length is limited in `2`,
  Only English letters are allowed. Changing this parameter will create a new resource.

* `state` - (Required, String, ForceNew) Specifies the state of private CA. The valid length is limited between
  `1` to `128`, Only Chinese and English letters, digits, hyphens (-), underscores (_), dots (.), comma (,) and
  space ( ) are allowed. Changing this parameter will create a new resource.

* `locality` - (Required, String, ForceNew) Specifies the locality of private CA. The valid length is limited between
  `1` to `128`, Only Chinese and English letters, digits, hyphens (-), underscores (_), dots (.), comma (,) and
  space ( ) are allowed. Changing this parameter will create a new resource.

* `organization` - (Required, String, ForceNew) Specifies the organization of private CA. The valid length is limited
  between `1` to `64`, Only Chinese and English letters, digits, hyphens (-), underscores (_), dots (.), comma (,) and
  space ( ) are allowed. Changing this parameter will create a new resource.

* `organizational_unit` - (Required, String, ForceNew) Specifies the organizational unit of private CA. The valid length
  is limited between `1` to `64`, Only Chinese and English letters, digits, hyphens (-), underscores (_), dots (.),
  comma (,) and space ( ) are allowed. Changing this parameter will create a new resource.

<a name="block-validity"></a>
The `validity` block supports:

* `type` - (Required, String, ForceNew) Specifies the type of validity value. Changing this parameter will create a new
  resource. Options are: **YEAR**, **MONTH(31 days)**, **DAY**, **HOUR**. If the charging mode is **prePaid**, only
  support **YEAR** and **MONTH(31 days)**.

* `value` - (Required, Int, ForceNew) Specifies the value of validity. Root CA certificate is no longer than 30 years
  and subordinate CA is no longer than 20 years. Changing this parameter will create a new resource. When creating a
  subordinate CA, the validity must less than the root CA.

* `started_at` - (Optional, String, ForceNew) Specifies the start time of validity. The value is a timestamp in milliseconds.
  For example, `1722840237000` indicates `2024-08-05 14:43:57`. The value of `started_at` cannot be earlier than `5` minutes
  from the current time.
  Changing this parameter will create a new resource.

<a name="block-crl_configuration"></a>
The `crl_configuration` block supports:

* `enabled` - (Optional, Bool) Specifies whether to enable the CRL publishing function. Defaults to **false**.

* `obs_bucket_name` - (Optional, String) Specifies the OBS bucket name. This field is required only when the CRL publishing
  function is enabled.

* `valid_days` - (Optional, Int) Specifies the CRL update interval, in days. It's limited from `7` to `30`.
  This field is required only when the CRL publishing function is enabled.

-> Due to API limitations, the fields `obs_bucket_name` and `valid_days` must be changed together with `enable`.
For example, if you want to update `obs_bucket_name` or `valid_days` when the CRL function is enabled, you need to
disable the CRL function first, then enable it and configure new `obs_bucket_name` or `valid_days` values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The private CA ID in UUID format.

* `status` - The current phase of the private CA. Valid values are as follows:
  + **PENDING**: The CA certificate is to be activated.
  + **ACTIVED**: The CA certificate is activated.
  + **DISABLED**: The CA certificate is disabled.
  + **DELETED**: The CA certificate is to be deleted as scheduled.
  + **EXPIRED**: The CA certificate has expired.

  Only private CAs with status `ACTIVED` can be used to issue certificates.

* `issuer_name` - The name of the parent CA. For a root CA, the value of this parameter is null.

* `gen_mode` - The generation method of the private CA. Valid values are as follows:
  + **GENERATE**: The certificate is generated through the PCA system.
  + **IMPORT**: The certificate is imported externally.
  + **CSR**: The CSR is imported externally and issued by the internal CA. The private key is not managed in PCA.

* `serial_number` - The serial number of the private CA.

* `created_at` - The creation time of the private CA.

* `expired_at` - The expiration time of the private CA.

* `free_quota` - The free quota of the private certificate.

* `crl_configuration/crl_dis_point` - The address of the CRL file in the OBS bucket.

* `crl_configuration/crl_name` - The name of the certificate revocation list.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Private CA can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ccm_private_ca.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `validity`, `key_usages`, `pending_days`,
`action`, `auto_renew`.
It is generally recommended running `terraform plan` after importing a private CA. You can then decide if changes should
be applied to it, also you can ignore changes as below.

```hcl
resource "huaweicloud_ccm_private_ca" "test" {
    ...

  lifecycle {
    ignore_changes = [
      validity, key_usages, pending_days, action, auto_renew,
    ]
  }
}
```
