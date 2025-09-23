---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_certificate"
description: |
  Manages a CCM private certificate resource within HuaweiCloud.
---

# huaweicloud_ccm_private_certificate

Manages a CCM private certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "common_name" {}
variable "issuer_id" {}

resource "huaweicloud_ccm_private_certificate" "test" {
  issuer_id           = var.issuer_id
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA256"
  
  distinguished_name {
    common_name = var.common_name
  }
  
  validity {
    type  = "DAY"
    value = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the certificate region. Changing this creates a new
  private certificate resource. Now only support cn-north-4 (china) and ap-southeast-3 (international)

* `distinguished_name` - (Required, List, ForceNew) Specifies the distinguished name of private certificate.
  Changing this parameter will create a new resource.
  The [distinguished_name](#block-distinguished_name) structure is documented below.

* `issuer_id` - (Required, String, ForceNew) Specifies the ID of the parent CA. Changing this creates
  a new private certificate resource.

* `key_algorithm` - (Required, String, ForceNew) Specifies the key algorithm. Valid values are **RSA2048**, **RSA4096**,
  **EC256**, and **EC384**.
  Changing this creates a new private certificate resource.

* `signature_algorithm` - (Required, String, ForceNew) Specifies the private certificate signature hash algorithm.
  Valid values are **SHA256**, **SHA384**, and **SHA512**. Changing this creates a new private certificate resource.

* `validity` - (Required, List, ForceNew) Specifies the validity of private certificate.
  Changing this parameter will create a new resource.
  The [validity](#block-validity) structure is documented below.

* `subject_alternative_names` - (Optional, List, ForceNew) Specifies the alternative name for the subject.
  Changing this parameter will create a new resource.
  The [subject_alternative_names](#block-subject_alternative_names) structure is documented below.

* `key_usage` - (Optional, List, ForceNew) Specifies the key usage. Changing this parameter will create a new resource.
  Refer to [4.2.1.3](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.3)
  Valid values are as follows:
  + **digitalSignature**: The key is used as a digital signature.
  + **nonRepudiation**: The key can be used for non-repudiation.
  + **keyEncipherment**: The key can be used for key encryption.
  + **dataEncipherment**: The key can be used for data encryption.
  + **keyAgreement**: The key can be used for key negotiation.
  + **keyCertSign**: The key can issue a certificate.
  + **cRLSign**: The key can issue a certificate revocation list (CRL).
  + **encipherOnly**: The key is used only for encryption.
  + **decipherOnly**: The key is used only for decryption.

* `server_auth` - (Optional, Bool, ForceNew) Specifies whether to enable server authentication. Defaults to **false**.
  Changing this parameter will create a new resource.

  -> Enable this enhanced key usage for the server certificate.

* `client_auth` - (Optional, Bool, ForceNew) Specifies whether to enable client authentication. Defaults to **false**.
  Changing this parameter will create a new resource.

  -> Enable this enhanced key usage for the client certificate.

* `code_signing` - (Optional, Bool, ForceNew) Specifies whether to enable signing of downloadable executable code client
  authentication. Defaults to **false**. Changing this parameter will create a new resource.

* `email_protection` - (Optional, Bool, ForceNew) Specifies whether to enable email protection. Defaults to **false**.
  Changing this parameter will create a new resource.

* `time_stamping` - (Optional, Bool, ForceNew) Specifies whether to bind the hash of an object to a time.
  Defaults to **false**. Changing this parameter will create a new resource.

* `object_identifier` - (Optional, String, ForceNew) Specifies the object identifier. Changing this parameter will create
  a new resource.
  
  -> The value of this parameter must be a dot-decimal notation string that complies with the ASN1 specifications,
  for example, **1.3.6.1.4.1.2011.4.99**.

* `object_identifier_value` - (Optional, String, ForceNew) Specifies the custom attribute content,
  work with `object_identifier`. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the certificate belongs.
  For enterprise users, if omitted, resources under all authorized enterprise projects of the tenant will be queried.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs associating with the private certificate.

<a name="block-distinguished_name"></a>
The `distinguished_name` block supports:

* `common_name` - (Required, String, ForceNew) Specifies the common name of private certificate. The valid length
  is limited between `1` to `64`, Only Chinese and English letters, digits, hyphens (-), underscores (_),
  dots (.), comma (,), space ( ) and asterisks (*) are allowed. Changing this parameter will create a new resource.

* `country` - (Optional, String, ForceNew) Specifies the country code of private certificate. The valid length is
  limited in `2`. Only English letters are allowed. The value corresponding to the parent CA is inherited by default.
  Changing this parameter will create a new resource.

* `state` - (Optional, String, ForceNew) Specifies the state or city name of private certificate. The valid length is
  limited between `1` to `128`. Only Chinese and English letters, digits, hyphens (-), underscores (_),
  dots (.), comma (,) and space ( ) are allowed. Changing this parameter will create a new resource.

* `locality` - (Optional, String, ForceNew) Specifies the locality of private certificate. The valid length
  is limited between `1` to `128`, Only Chinese and English letters, digits, hyphens (-), underscores (_),
  dots (.), comma (,) and space ( ) are allowed. The value corresponding to the parent CA is inherited by default.
  Changing this parameter will create a new resource.

* `organization` - (Optional, String, ForceNew) Specifies the organization of private certificate. The valid length
  is limited between `1` to `64`. Only Chinese and English letters, digits, hyphens (-), underscores (_), dots (.),
  comma (,) and space ( ) are allowed. The value corresponding to the parent CA is inherited by default.
  Changing this parameter will create a new resource.

* `organizational_unit` - (Optional, String, ForceNew) Specifies the organizational unit of private certificate.
  The valid length is limited between `1` to `64`. Only Chinese and English letters, digits, hyphens (-),
  underscores (_), dots (.), comma (,) and space ( ) are allowed. The value corresponding to the parent CA is inherited
  by default. Changing this parameter will create a new resource.

<a name="block-validity"></a>
The `validity` block supports:

* `type` - (Required, String, ForceNew) Specifies the validity period type. Changing this parameter will create a new
  resource. Options are: **YEAR**, **MONTH**(31 days), **DAY**, **HOUR**.

* `value` - (Required, Int, ForceNew) Specifies the value of validity period. The value of this parameter varies
  depending on the value of `type`. Changing this parameter will create a new resource.
  + Root CA certificates: no longer than `30` years.
  + Subordinate CA or private certificates: no longer than `20` years

* `start_at` - (Optional, String, ForceNew) Specifies the private certificate validity start from.
  The value is a timestamp in milliseconds. For example, `1645146939688` indicates `2022-02-18 09:15:39`.
  The value of `start_at` cannot be earlier than the result of the value of current time minus `5` minutes.
  Changing this creates a new private certificate resource.

<a name="block-subject_alternative_names"></a>
The `subject_alternative_names` block supports:

* `type` - (Required, String, ForceNew) Specifies the type of the alternative name. Currently,
  only **DNS**, **IP**, **EMAIL**, and **URI** are allowed. Changing this parameter will create a new resource.

* `value` - (Required, String, ForceNew) Specifies the value of the corresponding alternative name type.
  Changing this parameter will create a new resource.
  + When `type` is **DNS**, the value length ranges from `0` to `253` characters.  
  + When `type` is **IP**, the value length ranges from `0` to `39` characters. Support IPv4 and IPv6.  
  + When `type` is **EMAIL**, the value length ranges from `0` to `256` characters.  
  + When `type` is **URI**, the value length ranges from `0` to `253` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `issuer_name` - Indicates the parent CA name.

* `status` - Indicates the private certificate status. Valid values are: **ISSUED**, **EXPIRED** and **REVOKED**.

* `start_at` - Indicates the private certificate valid start time.

* `expired_at` - Indicates the private certificate valid expired time.

* `gen_mode` - Indicates the private certificate create mode,by system or user. Valid values are as follows:
  + **GENERATE**: The certificate is generated through the PCA system.
  + **IMPORT**: The certificate is imported externally.
  + **CSR**: The CSR is imported externally and issued by the internal CA. The private key is not managed in PCA.

* `created_at` - Indicates the time when the certificate was created.

## Import

The CCM private certificate can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ccm_private_certificate.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
  API response, security or some other reason. The missing attributes include: `validity`,`key_usage`,`server_auth`,
`client_auth`,`code_signing`,`email_protection`,`time_stamping`,`object_identifier`,`object_identifier_value`,
`subject_alternative_names`.
It is generally recommended running `terraform plan` after importing a private certificate. You can then decide if
changes should be applied to the private certificate, or the resource definition should be updated to align with the
private certificate. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_ccm_private_certificate" "test" {
  ...

  lifecycle {
    ignore_changes = [
        validity,
        key_usage,
        server_auth,
        client_auth,
        code_signing,
        email_protection,
        time_stamping,
        object_identifier,
        object_identifier_value,
        subject_alternative_names,
    ]
  }
}
```
