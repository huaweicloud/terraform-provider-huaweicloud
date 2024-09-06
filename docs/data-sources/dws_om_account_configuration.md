---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_om_account_configuration"
description: |-
  Use this data source to get the OM account configuration within HuaweiCloud.
---

# huaweicloud_dws_om_account_configuration

Use this data source to get the OM account configuration within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_om_account_configuration" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The current status of the OM account.

* `om_user_expires_time` - The expiration time of the OM account, in RFC3339 format.
