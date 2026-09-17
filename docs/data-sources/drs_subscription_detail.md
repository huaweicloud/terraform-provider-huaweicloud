---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_subscription_detail"
description: |-
  Use this data source to query the detail of a DRS subscription task within HuaweiCloud.
---

# huaweicloud_drs_subscription_detail

Use this data source to query the detail of a DRS subscription task within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_subscription_detail" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the subscription detail.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS subscription task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The name of the subscription task.

* `ip` - The intranet IP address.

* `enterprise_project_id` - The enterprise project ID.

* `status` - The status of the subscription task.

* `created_time` - The creation time of the subscription task, in timestamp format.

* `begin_time` - The start time of the subscription task, in timestamp format.

* `now_time` - The current time, in timestamp format.

* `engine_type` - The link type of the subscription task.

* `description` - The description of the subscription task.

* `subscription_data_type` - The subscribed data types.  
  The [subscription_data_type](#drs_subscription_detail_subscription_data_type) structure is documented below.

* `source_endpoint` - The source database instance information of the subscription.  
  The [source_endpoint](#drs_subscription_detail_source_endpoint) structure is documented below.

<a name="drs_subscription_detail_subscription_data_type"></a>
The `subscription_data_type` block supports:

* `is_dml_subscribed` - Whether DML is subscribed.

* `is_ddl_subscribed` - Whether DDL is subscribed.

<a name="drs_subscription_detail_source_endpoint"></a>
The `source_endpoint` block supports:

* `db_instance_id` - The database instance ID.

* `name` - The database name.

* `ip` - The database intranet IP address.

* `port` - The database port.

* `type` - The database type.

* `user_name` - The database user name.
