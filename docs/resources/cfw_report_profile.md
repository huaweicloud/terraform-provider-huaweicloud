---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_report_profile"
description: |-
  Manages a CFW report profile resource within HuaweiCloud.
---

# huaweicloud_cfw_report_profile

Manages a CFW report profile resource within HuaweiCloud.

## Example Usage

### Daily Report Profile

```hcl
variable "fw_instance_id" {}
variable "topic_urn" {}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = var.fw_instance_id
  category          = "daily"
  name              = "test-name"
  topic_urn         = var.topic_urn
  send_period       = "3"
  subscription_type = "0"
  status            = "0"
  description       = "test description"
}
```

### Weekly Report Profile

```hcl
variable "fw_instance_id" {}
variable "topic_urn" {}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = var.fw_instance_id
  category          = "weekly"
  name              = "test-name"
  topic_urn         = var.topic_urn
  send_period       = "1"
  send_week_day     = "1"
  subscription_type = "0"
  status            = "1"
  description       = "test description"
}
```

### Custom Report Profile

```hcl
variable "fw_instance_id" {}
variable "topic_urn" {}
variable "start_time" {
  type = number
}
variable "end_time" {
  type = number
}

resource "huaweicloud_cfw_report_profile" "test" {
  fw_instance_id    = var.fw_instance_id
  category          = "custom"
  name              = "test-name"
  topic_urn         = var.topic_urn
  subscription_type = "0"
  status            = "1"
  description       = "test description"

  statistic_period {
    start_time = var.start_time
    end_time   = var.end_time
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall ID.
  It is a unique ID generated after a firewall instance is created. You can obtain the firewall ID by referring to
  the data source `huaweicloud_cfw_firewalls`.

* `category` - (Required, String, NonUpdatable) Specifies the report type. The value can be:
  + **daily**: Daily report
  + **weekly**: Weekly report
  + **custom**: Custom report

* `name` - (Required, String) Specifies the template name.

* `topic_urn` - (Required, String) Specifies the topic URN.

* `send_period` - (Optional, String) Specifies the sending time. Valid value must be between `0` and `23`.
  This field is mandatory when `category` is set to **daily** or **weekly**.

* `send_week_day` - (Optional, String) Specifies the days in a week when data is sent.
  Valid value must be between `1` and `7`.
  This field is mandatory when `category` is set to **weekly**.

* `status` - (Optional, String) Specifies the enabling status. The value can be:
  + **0**: Disabled
  + **1**: Enabled

* `subscription_type` - (Optional, String) Specifies the notification method. The value can be:
  + **0**: SMN notification
  + **1**: No notification

* `statistic_period` - (Optional, List) Specifies the statistical period.
  This field is mandatory when `category` is set to **custom**.

  The [statistic_period](#statistic_period_struct) structure is documented below.

* `description` - (Optional, String) Specifies the description.

<a name="statistic_period_struct"></a>
The `statistic_period` block supports:

* `start_time` - (Optional, Int) Specifies the start time of the statistical period. Milliseconds-level timestamp.
  The value of this field must be the start of the day at midnight. Valid format examples: `1772035200000`.

* `end_time` - (Optional, Int) Specifies the end time of the statistical period. Milliseconds-level timestamp.
  The value of this field must be the end of the day at midnight. Valid format examples: `1772726399999`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the report profile ID).

* `topic_name` - The topic name.

## Import

The resource can be imported using `fw_instance_id`, `id` (report profile ID), separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_report_profile.test <fw_instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `description`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cfw_report_profile" "test" {
    ...

  lifecycle {
    ignore_changes = [
      description,
    ]
  }
}
```
