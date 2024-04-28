---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_certificate_export"
description: ""
---

# huaweicloud_ccm_private_certificate_export

Use this data source export a private Certificate within HuaweiCloud.

## Example Usage

```hcl
variable "certificate_id" {}

data "huaweicloud_ccm_private_certificate_export" "test3" {
  region         = "cn-north-4"
  type           = "OTHER"
  certificate_id = var.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the certificate region. Changing this creates a new
  private certificate resource. Now only support cn-north-4 (china) and ap-southeast-3 (international).

* `certificate_id` - (Required, String, ForceNew) Specifies the certificate ID of the private certificate
  you want to export.

* `type` - (Required, String, ForceNew) Specifies the Type of the server on which the certificate is installed.
  The options are as follows:
  + **APACHE**: This parameter is recommended if you want to use the certificate for an Apache server.
  + **NGINX**: This parameter is recommended if you want to use the certificate for an Nginx server.
  + **OTHER**: This parameter is recommended if you want to download a certificate in PEM format.
  + **IIS**: This parameter is recommended if you want to use the certificate for an IIS server.
  + **TOMCAT**: This parameter is recommended if you want to use the certificate for a TOMCAT server.

  -> **NOTE:** When the **type** is "IIS" or "TOMCAT" the export certificate file is different everytime.

* `sm_standard` - (Optional, String) Specifies the GB/T GMT0009 standard specifications and
  GB/T GMT0010 standard. When the certificate algorithm is SM2,
  it is only valid when it is passed in. If it is not passed in, it defaults to false.

* `password` - (Optional, String) Specifies the password used to encrypt the private key. Support the use of uppercase
  and lowercase English letters, numbers, special characters (such as.+- _ #), etc. The maximum length is 32 bytes.
  If not passed in, encryption will not be used for export by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `private_key` - Indicates the certificate private key in PEM format.

* `certificate` - Indicates the certificate content in PEM format.

* `certificate_chain` - Indicates the certificate chain in PEM format.

* `enc_certificate` - Indicates the encryption certificate content in PEM format.

* `enc_private_key` - Indicates the encryption certificate private key in PEM format.

* `enc_sm2_enveloped_key` - Indicates the National Security GMT0009 Standard Specification for Encrypting Private
  Keys SM2 Digital Envelope.
  
* `signed_and_enveloped_data` - Indicates the National Security GMT0010 Standard Specification for
  Encrypting Private Keys - Signature Digital Envelope.

* `keystore_pass` - Indicates the keystore password. If argument "password" passed in, it will be empty.

* `server_pfx` - Indicates the certificate file for IIS server. Encoding by base64.

* `server_jks` - Indicates the certificate file for TOMCAT server. Encoding by base64.
