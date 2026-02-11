---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_csrs"
description: |-
  Use this data source to get the list of CCM SSL CSR list.
---

# huaweicloud_ccm_csrs

Use this data source to get the list of CCM SSL CSR list.

## Example Usage

```hcl
data "huaweicloud_ccm_csrs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the CSR name.

* `private_key_algo` - (Optional, String) Specifies the key algorithm type. Valid values are:
  + **RSA_2048**
  + **RSA_3072**
  + **RSA_4096**
  + **EC_P256**
  + **EC_P384**
  + **SM2**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `csr_list` - The CSR list.

  The [csr_list](#csr_list_struct) structure is documented below.

<a name="csr_list_struct"></a>
The `csr_list` block supports:

* `id` - The CSR ID.

* `name` - The CSR name.

* `csr` - The CSR content.

* `domain_name` - The domain name bound to the CSR.

* `sans` - The additional domain name bound to the CSR.

* `private_key_algo` - The key algorithm.

* `usage` - The CSR usage.

* `company_country` - The country.

* `company_province` - The province.

* `company_city` - The city.

* `company_name` - The company name.

* `create_time` - The CSR creation time.

* `update_time` - The CSR update time.
