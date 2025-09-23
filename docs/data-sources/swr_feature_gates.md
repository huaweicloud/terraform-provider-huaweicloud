---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_feature_gates"
description: |-
  Use this data source to get the list of SWR feature gates.
---

# huaweicloud_swr_feature_gates

Use this data source to get the list of SWR feature gates.

## Example Usage

```hcl
data "huaweicloud_swr_feature_gates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `enable_experience` - Whether the experience center is enabled.

* `enable_hss_service` - Whether interconnection with HSS is enabled.

* `enable_image_scan` - Whether image scanning is enabled.

* `enable_sm3` - Whether SM algorithms are enabled.

* `enable_image_sync` - Whether image synchronization is enabled.

* `enable_cci_service` - Whether interconnection with CCI is enabled.

* `enable_image_label` - Whether image tagging is enabled.

* `enable_pipeline` - Whether pipeline is enabled.
