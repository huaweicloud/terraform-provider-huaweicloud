---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_email_send_records"
description: |-
  Use this data source to get the list of DAS email send records.
---

# huaweicloud_das_email_send_records

Use this data source to get the list of DAS email send records.

## Example Usage

```hcl
data "huaweicloud_das_email_send_records" "test" {
  datastore_type = "MySQL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the email send records are located.  
  If omitted, the provider-level region will be used.

* `datastore_type` - (Required, String) Specifies the database type.  
  The valid values are as follows:
  + **MySQL**
  + **TaurusDB**
  + **GaussDB**
  + **MariaDB**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of email send records that matched the filter parameters.  
  The [records](#email_send_records_attr) structure is documented below.

<a name="email_send_records_attr"></a>
The `records` block supports:

* `send_time` - The send time, in RFC3339 format.

* `status` - The send status.

* `email` - The email address.

* `topic_id` - The topic ID.

* `topic_urn` - The topic URN.

* `instance_health_reports` - The list of instance health reports.  
  The [instance_health_reports](#email_send_records_instance_health_reports_attr) structure is documented below.

<a name="email_send_records_instance_health_reports_attr"></a>
The `instance_health_reports` block supports:

* `task_id` - The report ID.

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

* `start_time` - The diagnosis start time, in RFC3339 format.

* `end_time` - The diagnosis end time, in RFC3339 format.
