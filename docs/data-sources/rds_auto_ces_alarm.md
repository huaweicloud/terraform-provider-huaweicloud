---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_auto_ces_alarm"
description: |-
  Use this data source to query the auto CES alarm configuration within HuaweiCloud.
---

# huaweicloud_rds_auto_ces_alarm

Use this data source to query the auto CES alarm configuration within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_rds_auto_ces_alarm" "test" {}
```

### Filter by Engine

```hcl
data "huaweicloud_rds_auto_ces_alarm" "test" {
  engine = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the auto CES alarm.  
  If omitted, the provider-level region will be used.

* `engine` - (Optional, String) Specifies the database engine to filter.  
  The valid values are as follows:
  + **mysql**: MySQL database engine.
  + **postgresql**: PostgreSQL database engine.
  + **sqlserver**: SQL Server database engine.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `entities` - The list of auto CES alarm configurations.  
  The [entities](#auto_ces_alarm_entities) structure is documented below.

<a name="auto_ces_alarm_entities"></a>
The `entities` block supports:

* `id` - The unique identifier of the alarm record.

* `domain_id` - The domain ID.

* `domain_name` - The domain name.

* `project_id` - The project ID.

* `project_name` - The project name.

* `engine_name` - The database engine name.

* `new_instance_default` - Whether to enable auto alarm for new instances by default.

* `switch_status` - The switch status of the auto alarm.

* `alarm_id` - The unique identifier of the alarm rule.

* `topic_urn` - The topic URN.

* `created_at` - The timestamp when the record was created.

* `updated_at` - The timestamp when the record was last updated.
