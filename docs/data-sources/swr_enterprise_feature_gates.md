---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_feature_gates"
description: |-
  Use this data source to get the global feature switch infos.
---

# huaweicloud_swr_enterprise_feature_gates

Use this data source to get the global feature switch infos.

## Example Usage

```hcl
data "huaweicloud_swr_enterprise_feature_gates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `enable_enterprise` - Indicates whether the enterprise project is enabled.

* `cer_available` - Indicates whether the SWR enterprise feature is enabled.

* `enable_user_def_obs` - Indicates whether the OBS bucket is enabled.
