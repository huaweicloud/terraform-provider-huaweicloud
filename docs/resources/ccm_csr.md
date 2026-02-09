---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_csr"
description: |-
  Manages a CCM SSL CSR resource within HuaweiCloud.
---

# huaweicloud_ccm_csr

Manages a CCM SSL CSR resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "domain_name" {}

resource "huaweicloud_ccm_csr" "test" {
  name             = var.name
  domain_name      = var.domain_name
  private_key_algo = "RSA_2048"
  usage            = "PERSONAL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the user-defined CSR name.

* `domain_name` - (Required, String, NonUpdatable) Specifies the domain name bound to the CSR. If you want to use the CSR,
  ensure the domain name bound to a certificate contains the domain name set here.

* `private_key_algo` - (Required, String, NonUpdatable) Specifies the private key algorithm. The value can be:
  + **RSA_2048**
  + **RSA_3072**
  + **RSA_4096**
  + **EC_P256**
  + **EC_P384**
  + **SM2**

* `usage` - (Required, String, NonUpdatable) Specifies the CSR usage. The value can be:
  + **PERSONAL**: Individual certificate.
  + **ENTERPRISE**: Enterprise certificate.

* `sans` - (Optional, String, NonUpdatable) Specifies the additional domain name bound to the CSR.

* `company_country` - (Optional, String, NonUpdatable) Specifies the country.
  This parameter is mandatory when `usage` is set to **ENTERPRISE**. Example value: **CN**.

* `company_province` - (Optional, String, NonUpdatable) Specifies the province.
  This parameter is mandatory when `usage` is set to **ENTERPRISE**. Example value: **Beijing**.

* `company_city` - (Optional, String, NonUpdatable) Specifies the city.
  This parameter is mandatory when `usage` is set to **ENTERPRISE**. Example value: **Beijing**.

* `company_name` - (Optional, String, NonUpdatable) Specifies the company name.
  This parameter is mandatory when `usage` is set to **ENTERPRISE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the CSR ID.

* `csr` - The CSR content.

* `create_time` - The creation time of the CSR.

* `update_time` - The update time of the CSR.

## Import

The CCM SSL CSR can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ccm_csr.test <id>
```
