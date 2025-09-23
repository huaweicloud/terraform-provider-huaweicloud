---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_om_account_action"
description: |-
  Use this resource to operate OM account for DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_om_account_action

Use this resource to operate OM account for DWS cluster within HuaweiCloud.

-> 1. This resource is supported only in `8.1.3.110` or later.
   <br>2. This resource is only a one-time action resource for operating the OM account. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Turn on the OM account switch

```hcl
variable "dws_cluster_id" {}

resource "huaweicloud_dws_om_account_action" "test" {
  cluster_id = var.dws_cluster_id
  operation  = "addOmUser"
}
```

### The validity period of the OM account is extended by 16 hours

```hcl
variable "dws_cluster_id" {}

resource "huaweicloud_dws_om_account_action" "test" {
  cluster_id = var.dws_cluster_id
  operation  = "addOmUser"
}

resource "huaweicloud_dws_om_account_action" "increase_period" {
  count      = 2
  depends_on = [huaweicloud_dws_om_account_action.test]
  cluster_id = var.dws_cluster_id
  operation  = "increaseOmUserPeriod"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID.
  Changing this creates a new resource.

* `operation` - (Required, String, ForceNew) Specifies the operation type of the OM account.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **addOmUser**: Turn on the OM account switch.
  + **deleteOmUser**: Turn off the OM account switch.
  + **increaseOmUserPeriod**: Extend the validity period of the OM account.

  -> The **increaseOmUserPeriod** action is available only the OM account is enabled.
     Each time a resource is created, the account validity period is extended by `8` hours
     based on the current expiration time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
