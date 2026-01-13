---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_alarm_config"
description: |-
  Use this data source to get the Advanced Anti-DDos alarm config.
---

# huaweicloud_aad_alarm_config

Use this data source to get the Advanced Anti-DDos alarm config.

## Example Usage

```hcl
data "huaweicloud_aad_alarm_config" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `blackhole` - The blackhole.

* `ddos` - The DDOS.

* `topic_name` - The topic name.

* `topic_urn` - The topic URN.
