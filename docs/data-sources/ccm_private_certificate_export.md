---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_certificate_export"
description: |
  Use this data source export a private Certificate within HuaweiCloud.
---

# huaweicloud_ccm_private_certificate_export

Use this data source export a private Certificate within HuaweiCloud.

## Example Usage

```hcl
variable "certificate_id" {}

data "huaweicloud_ccm_private_certificate_export" "test" {
  region         = "cn-north-4"
  type           = "OTHER"
  certificate_id = var.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the certificate region in which to query the resource.
  If omitted, the provider-level region will be used.

* `certificate_id` - (Required, String) Specifies the certificate ID of the private certificate
  you want to export.

* `type` - (Required, String) Specifies the type of the server on which the certificate is installed.
  The options are as follows:
  + **APACHE**: Using for apache server.
  + **NGINX**: Using for nginx server.
  + **OTHER**: Using for download certificates in PEM format.
  + **IIS**: Using for Windows server.
  + **TOMCAT**: Using for tomcat server.

  -> The certificate file exported is different each time when `type` is set to **IIS** or **TOMCAT**.

* `sm_standard` - (Optional, String) Specifies whether to use the national secret **GMT0009** and **GMT0010** standard
  specification. This field is valid only when the certificate algorithm is **SM2**.
  The sm2 cert only support **OTHER** type. Valid values are **true** and **false**. Defaults to **false**.

* `password` - (Optional, String) Specifies the password used to encrypt the private key. Only uppercase letters,
  lowercase letters, digits, and special characters (`,.+-_#`) are allowed. The maximum length is 32 bytes.
  By default, encryption is not used when exporting.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `private_key` - Indicates the certificate private key in PEM format.

* `certificate` - Indicates the certificate content in PEM format.

* `certificate_chain` - Indicates the certificate chain in PEM format.

* `enc_certificate` - Indicates the encryption certificate content in PEM format.

* `enc_private_key` - Indicates the encryption certificate private key in PEM format.

* `enc_sm2_enveloped_key` - Indicates the national secret **GMT0009** standard specification SM2 digital envelope for
  encrypting private key.

* `signed_and_enveloped_data` - Indicates the national secret **GMT0010** standard specification signed digital envelope
  with encrypted private key.

* `keystore_pass` - Indicates the keystore password. This field is empty when argument `password` is specified.

* `server_pfx` - Indicates the certificate file for **IIS** server. Encoding by base64.

* `server_jks` - Indicates the certificate file for **TOMCAT** server. Encoding by base64.
