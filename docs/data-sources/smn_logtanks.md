---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_logtanks"
description: |-
  Use this data source to get a list of SMN logtanks.
---

# huaweicloud_smn_logtanks

Use this data source to get a list of SMN logtanks.

## Example Usage

```hcl
variable "topic_urn" {}

data "huaweicloud_smn_logtanks" "test" {
  topic_urn = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `topic_urn` - (Required, String) Specifies the topic URN.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logtanks` - The cloud logs.

  The [logtanks](#logtanks_struct) structure is documented below.

<a name="logtanks_struct"></a>
The `logtanks` block supports:

* `id` - The ID of the cloud log.

* `log_group_id` - The LTS log group ID.

* `log_stream_id` - The LTS log stream ID.

* `created_at` - The creation time.

* `updated_at` - The update time.
