---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca_agency"
description: |-
  Use this data source to get the list of CCM private CA agency.
---

# huaweicloud_ccm_private_ca_agency

Use this data source to get the list of CCM private CA agency.

## Example Usage

```hcl
data "huaweicloud_ccm_private_ca_agency" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agency_granted` - The agency granted.
