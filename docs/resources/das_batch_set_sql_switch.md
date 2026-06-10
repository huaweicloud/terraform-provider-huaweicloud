---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_batch_set_sql_switch"
description: |-
  Use this resource to batch set full SQL or slow SQL switch within HuaweiCloud.
---

# huaweicloud_das_batch_set_sql_switch

Use this resource to batch set full SQL or slow SQL switch within HuaweiCloud.

-> This resource is a one-time action resource for batch setting SQL switch. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Batch disable slow SQL switch of MySQL instances

```hcl
variable "instance_ids" {
  type = list(string)
}

resource "huaweicloud_das_batch_set_sql_switch" "test" {
  instance_ids    = var.instance_ids
  engine_type     = "MySQL"
  switch_on       = false
  switch_type     = "slowsql"
  retention_hours = 168
}
```

### Batch enable full SQL switch of SqlServer instances

```hcl
variable "instance_ids" {
  type = list(string)
}

resource "huaweicloud_das_batch_set_sql_switch" "test" {
  instance_ids    = var.instance_ids
  engine_type     = "SqlServer"
  switch_on       = true
  switch_type     = "fullsql"
  retention_hours = 300
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SQL switch is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_ids` - (Required, List, NonUpdatable) Specifies the list of instance IDs.

* `engine_type` - (Required, String, NonUpdatable) Specifies the engine type of the instances.  
  The valid values are as follows:
  + **MySQL**
  + **SqlServer**
  + **Taurus**
  + **MariaDB**

* `switch_on` - (Required, Bool, NonUpdatable) Whether to enable the SQL switch.  
  The valid values are as follows:
  + **true**: Enable the SQL switch.
  + **false**: Disable the SQL switch.

* `switch_type` - (Required, String, NonUpdatable) Specifies the type of SQL switch to set.  
  The valid values are as follows:
  + **fullsql**: The full SQL switch.
  + **slowsql**: The slow SQL switch.

* `retention_hours` - (Optional, Int, NonUpdatable) Specifies the retention hours of the SQL data.  
  + For **slow sql**, the valid value is range from `24` to `720`. Defaults to `168`.
  + For **full sql**, the valid value is range from `24` to `4320`. Defaults to `168`.

  -> The `retention_hours` is only valid when `switch_on` is set to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
