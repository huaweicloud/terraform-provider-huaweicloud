---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_event_overview"
description: |-
  Use this data source to get the statistical overview of security events within HuaweiCloud.
---

# huaweicloud_dsc_event_overview

Use this data source to get the statistical overview of security events within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_event_overview" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `block_num` - The number of blocking events.

* `turn_off_num` - The number of events in the off state.

* `turn_on_num` - The number of events in the on state.

* `not_overdue_event` - The statistical information of non-overdue events.

  The [not_overdue_event](#not_overdue_event_struct) structure is documented below.

* `overdue_event` - The statistical information of overdue events.

  The [overdue_event](#overdue_event_struct) structure is documented below.

<a name="not_overdue_event_struct"></a>
The `not_overdue_event` block supports:

* `fatal_num` - The number of fatal events.

* `high_risk_num` - The number of high-risk events.

* `middle_risk_num` - The number of medium-risk events.

* `low_risk_num` - The number of low-risk events.

* `notice_num` - The number of notice events.

<a name="overdue_event_struct"></a>
The `overdue_event` block supports:

* `fatal_num` - The number of fatal events.

* `high_risk_num` - The number of high-risk events.

* `middle_risk_num` - The number of medium-risk events.

* `low_risk_num` - The number of low-risk events.

* `notice_num` - The number of notice events.
