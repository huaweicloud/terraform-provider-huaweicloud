---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_port_info"
description: |-
  Use this data source to get the detail information of a specifies port.
---

# huaweicloud_hss_asset_port_info

Use this data source to get the detailed information of a specifies port.

## Example Usage

```hcl
data "huaweicloud_hss_asset_port_info" "test" {
  port     = 8080
  category = "0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `port` - (Required, Int) Specifies the port number.

* `category` - (Required, String) Specifies the asset type.
  The valid values are as follows:
  + **0**: Host.
  + **1**: Container.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `type` - The port type.

* `status` - The status of the port.
  The valid values are as follows:
  + **normal**: Normal.
  + **danger**: Dangerous.
  + **unknow**: Unknown.

* `description` - The description of the port in Chinese.

* `description_en` - The description of the port in English.
