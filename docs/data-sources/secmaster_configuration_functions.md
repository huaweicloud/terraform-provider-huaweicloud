---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_configuration_functions"
description: |-
  Use this data source to get the list of SecMaster configuration functions within HuaweiCloud.
---

# huaweicloud_secmaster_configuration_functions

Use this data source to get the list of SecMaster configuration functions within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_secmaster_configuration_functions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `support_postpaid_mgmt` - Whether the pay-per-use mode is supported.
  + **true**: Supported.
  + **false** or not returned: Not supported.

* `support_large_screen_mgmt` - Whether the large screen management is supported.
  + **true**: Supported.
  + **false** or not returned: Not supported.

* `support_purchase_label_mgmt` - Whether the label management is supported.
  + **true**: Supported.
  + **false** or not returned: Not supported.

* `billing_type_mgmt` - The billing policy.
  The valid values are as follows:
  + **CBC**
  + **BSS**
