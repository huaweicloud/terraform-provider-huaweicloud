---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_associate_smn_topic"
description: |-
  Use this resource to associate an SMN topic with a DRS job within HuaweiCloud.
---

# huaweicloud_drs_associate_smn_topic

Use this resource to associate an SMN topic with a DRS job within HuaweiCloud.

-> This resource is a one-time action resource for associating an SMN topic with a DRS job. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}
variable "topic_urn" {}

resource "huaweicloud_drs_associate_smn_topic" "test" {
  job_id                    = var.job_id
  topic_urn                 = var.topic_urn
  is_alarm_to_user          = true
  increment_delay_threshold = 500
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to associate the SMN topic.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS job to associate the SMN topic.

* `topic_urn` - (Required, String, NonUpdatable) Specifies the topic URN of the SMN topic.

* `is_alarm_to_user` - (Required, Bool, NonUpdatable) Specifies whether to send alarm notifications to the user.

* `increment_delay_threshold` - (Optional, Int, NonUpdatable) Specifies the incremental delay threshold in seconds.

* `rto_delay_threshold` - (Optional, Int, NonUpdatable) Specifies the RTO delay threshold in seconds.  
  This parameter is only available in disaster recovery scenarios.

* `rpo_delay_threshold` - (Optional, Int, NonUpdatable) Specifies the RPO delay threshold in seconds.  
  This parameter is only available in disaster recovery scenarios.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the job ID.
