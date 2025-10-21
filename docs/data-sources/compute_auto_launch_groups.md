---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_auto_launch_groups"
description: |-
  Use this data source to get the list of auto launch groups.
---

# huaweicloud_compute_auto_launch_groups

Use this data source to get the list of auto launch groups.

## Example Usage

```hcl
data "huaweicloud_compute_auto_launch_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the auto launch group.

* `valid_since` - (Optional, String) Specifies the request start time.
  The value is in the format of **yyyy-MM-ddTHH:mm:ssZ** in UTC+0 and complies with ISO8601.

* `valid_until` - (Optional, String) Specifies the request end time.
  The value is in the format of **yyyy-MM-ddTHH:mm:ssZ** in UTC+0 and complies with ISO8601.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `auto_launch_groups` - Indicates the list of auto launch groups.

  The [auto_launch_groups](#auto_launch_groups_struct) structure is documented below.

<a name="auto_launch_groups_struct"></a>
The `auto_launch_groups` block supports:

* `id` - Indicates the ID of the auto launch group.

* `name` - Indicates the name of the auto launch group.

* `type` - Indicates the request type.

* `status` - Indicates the status of the auto launch group.
  The value can be: **SUBMITTED**, **ACTIVE**, **DELETING**, **DELETED**.

* `task_state` - Indicates the status of the auto launch group task.
  The value can be:
  + **HANDLING**: Launching.
  + **FULFILLED**: The auto launch group task is fully equipped.
  + **ERROR**: Error occurs in the auto launch group task.

* `valid_since` - Indicates the request start time.

* `valid_until` - Indicates the request end time.
