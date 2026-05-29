---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_history_transactions"
description: |-
  Use this data source to get the list of DAS history transactions.
---

# huaweicloud_das_history_transactions

Use this data source to get the list of DAS history transactions.

## Basic Usage

### Query transactions by order

```hcl
variable "instance_id" {}

data "huaweicloud_das_history_transactions" "test" {
  instance_id    = var.instance_id
  datastore_type = "MySQL"
  start_time     = "2026-06-01T00:00:00+08:00"
  end_time       = "2026-06-06T00:00:00+08:00"
  order_by       = "desc"
  order_field    = "lastSec"
}
```

### Query transactions by order and duration

```hcl
variable "instance_id" {}

data "huaweicloud_das_history_transactions" "test" {
  instance_id    = var.instance_id
  datastore_type = "MySQL"
  start_time     = "2026-06-01T00:00:00+08:00"
  end_time       = "2026-06-06T00:00:00+08:00"
  last_sec_min   = 1
  last_sec_max   = 60
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the history transactions are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the database instance.

* `datastore_type` - (Required, String) Specifies the database type. The valid value is **MySQL**.

* `start_time` - (Required, String) Specifies the start time of the query range, in RFC3339 format.

* `end_time` - (Required, String) Specifies the end time of the query range, in RFC3339 format.

* `order_by` - (Optional, String) Specifies the sort order.  
  The valid values are as follows:
  + **asc**
  + **desc**

* `order_field` - (Optional, String) Specifies the field used for sorting.  
  The valid values are as follows:
  + **lastSec**. The duration.
  + **waitLockStructCount**. The number of wait locks.
  + **holdLockStructCount**. The number of hold locks.
  + **collectTime**. The collection time.

* `last_sec_min` - (Optional, Int) Specifies the minimum duration of the transaction, in seconds.

* `last_sec_max` - (Optional, Int) Specifies the maximum duration of the transaction, in seconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `transactions` - The list of history transactions that matched the filter parameters.  
  The [transactions](#history_transactions_attr) structure is documented below.

<a name="history_transactions_attr"></a>
The `transactions` block supports:

* `last_sec` - The transaction duration, in seconds.

* `wait_locks` - The number of wait locks.

* `hold_locks` - The number of hold locks.

* `occurrence_time` - The occurrence time, in RFC3339 format.

* `detail` - The transaction content.

* `collect_time` - The collect time, in RFC3339 format.
