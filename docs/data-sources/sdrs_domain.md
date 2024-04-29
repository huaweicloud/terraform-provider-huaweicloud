---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_domain"
description: ""
---

# huaweicloud_sdrs_domain

Use this data source to get an available SDRS domain.

## Example Usage

```hcl
data "huaweicloud_sdrs_domain" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of an available SDRS domain.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `description` - Indicates the description of the SDRS domain.
