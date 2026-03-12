---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_monitor_indicators"
description: |-
  Use this data source to query monitor indicators of DWS within HuaweiCloud.
---

# huaweicloud_dws_monitor_indicators

Use this data source to query monitor indicators of DWS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_monitor_indicators" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region where the monitor indicators are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `indicators` - The list of monitor indicators.  
  The [indicators](#dws_indicators_struct) structure is documented below.

<a name="dws_indicators_struct"></a>
The `indicators` block supports:

* `indicator_name` - The monitor indicator name.

* `plugin_name` - The collection plugin name.

* `default_collect_rate` - The default collection rate.

* `support_datastore_version` - The supported datastore version.
