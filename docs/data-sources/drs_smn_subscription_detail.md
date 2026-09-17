---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_smn_subscription_detail"
description: |-
  Use this data source to query the SMN subscription detail of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_smn_subscription_detail

Use this data source to query the SMN subscription detail of a DRS job within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_smn_subscription_detail" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the SMN subscription detail.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job to query SMN subscription detail.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topic_name` - The name of the SMN topic.

* `topic_urn` - The unique resource identifier of the SMN topic.

* `delay_time` - The delay time in seconds.

* `rpo_delay` - The RPO delay time in seconds.

* `rto_delay` - The RTO delay time in seconds.

* `is_alarm_user` - Whether to send alarm to the subscribed user.
