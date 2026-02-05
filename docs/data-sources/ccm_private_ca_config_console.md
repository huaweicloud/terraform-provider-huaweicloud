---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca_config_console"
description: |-
  Use this data source to get the CCM private CA config console.
---

# huaweicloud_ccm_private_ca_config_console

Use this data source to get the CCM private CA config console.

## Example Usage

```hcl
data "huaweicloud_ccm_private_ca_config_console" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_support_dhsm` - Whether the Dedicated HSM cluster is supported.

* `dhsm_regions` - Regions supported by the Dedicated HSM cluster.

* `is_support_yearly_monthly_ca` - Whether yearly/monthly CAs are supported.

* `is_support_iam5` - Whether IAM5 authentication is supported.

* `is_support_ocsp` - Whether OCSP queries are supported.

* `is_support_eps` - Whether enterprise projects are supported.

* `is_support_sm2` - Whether the SM2 algorithm is supported.
