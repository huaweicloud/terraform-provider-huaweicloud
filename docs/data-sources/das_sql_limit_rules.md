---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_sql_limit_rules"
description: |-
  Use this data source to query DAS SQL limit rules within HuaweiCloud.
---

# huaweicloud_das_sql_limit_rules

Use this data source to query DAS SQL limit rules within HuaweiCloud.

-> This data source only supports to query SQL limit rules of **MySQL** instances.

## Example Usage

### Query all SQL limit rules under a specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_das_sql_limit_rules" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the SQL limit rules are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the instance to which the SQL limit rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `rules` - The list of SQL limit rules that matched filter parameters.  
  The [rules](#das_sql_limit_rules) structure is documented below.

<a name="das_sql_limit_rules"></a>
The `rules` block supports:

* `id` - The ID of the SQL limit rule.

* `sql_type` - The type of the SQL.

* `pattern` - The pattern of the SQL limit rule.

* `max_concurrency` - The maximum concurrency of the SQL limit rule.

* `max_waiting` - The maximum waiting time of the SQL limit rule.
