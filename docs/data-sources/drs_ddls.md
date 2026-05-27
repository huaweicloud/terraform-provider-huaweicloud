---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_ddls"
description: |-
  Use this data source to get the list of DDL alarm information of a DRS task within HuaweiCloud.
---

# huaweicloud_drs_ddls

Use this data source to get the list of DDL alarm information of a DRS task within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_ddls" "test" {
  job_id = "your_job_id"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS task.

* `start_seq_no` - (Optional, String) Specifies the start value of the DDL sequence.

* `end_seq_no` - (Optional, String) Specifies the end value of the DDL sequence.

* `status` - (Optional, String) Specifies the DDL status.
  The valid values are as follows:
  + **0**: No alarm.
  + **1**: Alarm in progress.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ddls` - The list of DDL alarm information.

  The [ddls](#ddls_struct) structure is documented below.

<a name="ddls_struct"></a>
The `ddls` block supports:

* `seqno` - The unique sequence number of the record.

* `checkpoint` - The checkpoint of the source database.

* `status` - The DDL alarm status.
  + **0**: No alarm.
  + **1**: Alarm in progress.

* `ddl_timestamp` - The time when the DDL is executed in the source database.

* `ddl_text` - The content of the DDL.

* `exe_result` - The execution result of the DDL.
  + **false**: Execution failed.
  + **true**: Execution succeeded.

* `record_time` - The time when the data is recorded.

* `clean_time` - The time when the DDL alarm is cleared.
