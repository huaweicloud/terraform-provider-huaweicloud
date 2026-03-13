---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ip_blacklist_switch"
description: |-
  Use this data source to get the traffic filtering IP blacklist switch information within HuaweiCloud.
---

# huaweicloud_cfw_ip_blacklist_switch

Use this data source to get the traffic filtering IP blacklist switch information within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_ip_blacklist_switch" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The IP blacklist switch information.

  The [data](#cfw_ip_blacklist_switch_data) structure is documented below.

<a name="cfw_ip_blacklist_switch_data"></a>
The `data` block supports:

* `status` - The IP blacklist function switch status.  
  The valid values are as follows:
  + **1**: Enable.
  + **0**: Disable.
