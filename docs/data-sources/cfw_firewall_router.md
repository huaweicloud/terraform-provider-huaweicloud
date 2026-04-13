---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_firewall_router"
description: |-
  Use this data source to get the enterprise router list associated with the CFW east-west firewall within HuaweiCloud.
---

# huaweicloud_cfw_firewall_router

Use this data source to get the enterprise router list associated with the CFW east-west firewall within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_firewall_router" "test" { 
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

* `id` - The data source ID.

* `er_list` - The list of enterprise routers associated with the east-west firewall.

  The [er_list](#er_list_struct) structure is documented below.

<a name="er_list_struct"></a>
The `er_list` block supports:

* `er_id` - The enterprise router ID.

* `name` - The enterprise router name.
