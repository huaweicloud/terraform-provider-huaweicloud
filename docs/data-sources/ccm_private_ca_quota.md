---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca_quota"
description: |-
  Use this data source to get the quota of CCM private CA.
---

# huaweicloud_ccm_private_ca_quota

Use this data source to get the quota of CCM private CA.

## Example Usage

```hcl
data "huaweicloud_ccm_private_ca_quota" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The certificate quota.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - The resource quota list.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - The certificate type. Valid values are:
  + **CERTIFICATE_AUTHORITY**: CA certificate.
  + **CERTIFICATE**: Private certificate.

* `used` - The used quota.

* `quota` - The total quota. There are two scenarios:
  + For **CERTIFICATE_AUTHORITY** certificate, the default value is `100`.
  + For **CERTIFICATE** certificate, the default value is `100,000`.
