---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_event_logs"
description: |-
  Use this data source to get the event logs of WAF within HuaweiCloud.
---

# huaweicloud_waf_event_logs

Use this data source to get the event logs of WAF within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_waf_event_logs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - List of files to download for incident protection.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `description` - The event description.

* `filename` - The filename.

* `obsname` - The file OBS name.

* `start` - The statistics start time.

* `source` - The file source.

* `state` - The file state.

* `id` - The file ID.

* `end` - The statistics deadline.

* `url` - The URL.

* `urltimestamp` - The update URL timestamp.

* `timestamp` - The file generation timestamp.
