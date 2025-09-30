---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_maintainwindows"
description: |-
  Use this data source to get the list of Kafka maintain windows within HuaweiCloud.
---

# huaweicloud_dms_kafka_maintainwindows

Use this data source to get the list of Kafka maintain windows within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dms_kafka_maintainwindows" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the maintain windows are located.  
  If omitted, the provider-level region will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `maintain_windows` - The list of the maintain windows.  
  The [maintain_windows](#kafka_maintain_windows_attribute) attribute is documented below.

<a name="kafka_maintain_windows_attribute"></a>
The `maintain_windows` block supports:

* `default` - Whether this is the default time window.

* `begin` - The start time of the maintain window.

* `end` - The end time of the maintain window.

* `seq` - The sequence number.
