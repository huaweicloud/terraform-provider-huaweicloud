---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_sql_limit_rule"
description: |-
  Manages a DAS SQL limit rule resource within HuaweiCloud.
---

# huaweicloud_das_sql_limit_rule

Manages a DAS SQL limit rule resource within HuaweiCloud.

-> This resource only supports to manage the SQL limit rules of **MySQL** instances.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_das_sql_limit_rule" "test" {
  instance_id     = var.instance_id
  sql_type        = "SELECT"
  pattern         = "select~test"
  max_concurrency = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SQL limit rule is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the database instance.

* `sql_type` - (Required, String, NonUpdatable) Specifies the SQL type.  
  The valid values are as follows:
  + **SELECT**
  + **UPDATE**
  + **DELETE**

* `pattern` - (Required, String, NonUpdatable) Specifies the SQL limit rule pattern.  
  For example, the keyword `select~test` means that **select** and **test** are the two keywords contained in this
  concurrency control, and **~** is the keyword separator. If the SQL command contains the two keywords
  **select** and **test**, it is considered to have successfully matched this concurrency control rule.

* `max_concurrency` - (Required, Int, NonUpdatable) Specifies the maximum concurrency.  
  The valid value is range from **0** to **2^31-1**.

* `database_name` - (Optional, String, NonUpdatable) Specifies the database name.

* `max_waiting` - (Optional, Int, NonUpdatable) Specifies the maximum waiting time.  
  The valid value is range from **0** to **2^31-1**.

* `his_sql_limit_switch` - (Optional, Bool, NonUpdatable) Whether to enable the historical SQL limit switch.  
  The existing sessions that match this rule will be killed when the parameter is set to **true**.  
  The valid values are as follows:
  + **true**
  + **false**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

## Import

The SQL limit rule can be imported using `<instance_id>/<rule_id>`, e.g.

```bash
$ terraform import huaweicloud_das_sql_limit_rule.test <instance_id>/<rule_id>
```
