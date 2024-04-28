---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_recorder"
description: ""
---

# huaweicloud_rms_resource_recorder

Manages a RMS recorder resource within HuaweiCloud.

-> Only one resource recorder can be configured.

## Example Usage

### Recorder with All Supported Resources

```hcl
variable "topic_urn" {}
variable "bucket_name" {}
variable "delivery_region" {}

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported = true
  }

  obs_channel {
    bucket = var.bucket_name
    region = var.delivery_region
  }
  smn_channel {
    topic_urn = var.topic_urn
  }
}
```

### Recorder with Specified Resources

```hcl
variable "bucket_name" {}
variable "delivery_region" {}

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported  = false
    resource_types = ["vpc.vpcs", "rds.instances", "dms.kafkas", "dms.rabbitmqs", "dms.queues"]
  }

  obs_channel {
    bucket = var.bucket_name
    region = var.delivery_region
  }
}
```

## Argument Reference

The following arguments are supported:

* `agency_name` - (Required, String) Specifies the IAM agency name which must include permissions
  for sending notifications through SMN and for writing data into OBS.

* `selector` - (Required, List) Specifies configurations of resource selector.
  The [object](#Recorder_SelectorConfigBody) structure is documented below.

* `obs_channel` - (Optional, List) Specifies configurations of the OBS bucket used for data dumping.
  The [object](#Recorder_TrackerOBSChannelConfigBody) structure is documented below.

* `smn_channel` - (Optional, List) Specifies configurations of the SMN channel used to send notifications.
  The [object](#Recorder_TrackerSMNChannelConfigBody) structure is documented below.

-> At least one `obs_channel` or `smn_channel` must be configured.

<a name="Recorder_SelectorConfigBody"></a>
The `selector` block supports:

* `all_supported` - (Required, Bool) Specifies whether to select all supported resources.

* `resource_types` - (Optional, List) Specifies the resource type list. For details about the supported services and
  resource types, see [Supported Services and Resource Types](https://support.huaweicloud.com/intl/en-us/api-rms/rms_06_0100.html).

<a name="Recorder_TrackerOBSChannelConfigBody"></a>
The `obs_channel` block supports:

* `bucket` - (Required, String) Specifies the OBS bucket name used for data dumping.

* `region` - (Required, String) Specifies the region where this bucket is located.

<a name="Recorder_TrackerSMNChannelConfigBody"></a>
The `smn_channel` block supports:

* `topic_urn` - (Required, String) Specifies the SMN topic URN used to send notifications.

* `region` - (Optional, String) Specifies the region where this SMN topic is located.

* `project_id` - (Optional, String) Specifies the project ID where this SMN topic is located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The recorder can be imported using **any string** as the resource ID, e.g.

```bash
$ terraform import huaweicloud_rms_resource_recorder.test recorder
```
