---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_status"
description: |-
  Use this data source to query the number of DCS instances in different statuses within HuaweiCloud.
---

# huaweicloud_dcs_instance_status

Use this data source to query the number of DCS instances in different statuses within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dcs_instance_status" "test" {
  include_failure = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the status. If omitted, the
  provider-level region will be used.

* `include_failure` - (Optional, String) Whether to include the number of instances that failed to be created.
  The value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `paying_count` - The number of instances in the **PAYING** status.

* `freezing_count` - The number of instances in the **FREEZING** status.

* `migrating_count` - The number of instances in the **MIGRATING** status.

* `flushing_count` - The number of instances in the **FLUSHING** status.

* `upgrading_count` - The number of instances in the **UPGRADING** status.

* `restoring_count` - The number of instances in the **RESTORING** status.

* `extending_count` - The number of instances in the **EXTENDING** status.

* `creating_count` - The number of instances in the **CREATING** status.

* `running_count` - The number of instances in the **RUNNING** status.

* `error_count` - The number of instances in the **ERROR** status.

* `frozen_count` - The number of instances in the **FROZEN** status.

* `createfailed_count` - The number of instances that failed to be created.

* `restarting_count` - The number of instances in the **RESTARTING** status.

* `redis` - The status statistics of Redis instances.
  The [redis](#status_statistic_struct) structure is documented below.

* `memcached` - The status statistics of Memcached instances.
  The [memcached](#status_statistic_struct) structure is documented below.

<a name="status_statistic_struct"></a>
The `redis` and `memcached` block supports:

* `paying_count` - The number of instances in the **PAYING** status.

* `freezing_count` - The number of instances in the **FREEZING** status.

* `migrating_count` - The number of instances in the **MIGRATING** status.

* `flushing_count` - The number of instances in the **FLUSHING** status.

* `upgrading_count` - The number of instances in the **UPGRADING** status.

* `restoring_count` - The number of instances in the **RESTORING** status.

* `extending_count` - The number of instances in the **EXTENDING** status.

* `creating_count` - The number of instances in the **CREATING** status.

* `running_count` - The number of instances in the **RUNNING** status.

* `error_count` - The number of instances in the **ERROR** status.

* `frozen_count` - The number of instances in the **FROZEN** status.

* `createfailed_count` - The number of instances that failed to be created.

* `restarting_count` - The number of instances in the **RESTARTING** status.
