---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_sql_limiting_switch"
description: |-
  Use this resource to enable or disable the SQL limiting switch within HuaweiCloud.
---

# huaweicloud_das_sql_limiting_switch

Use this resource to enable or disable the SQL limiting switch within HuaweiCloud.

-> This resource is a one-time action resource for switching the SQL limiting. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_das_sql_limiting_switch" "test" {
  instance_id    = var.instance_id
  datastore_type = "MySQL"
  status         = "ON"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SQL limiting switch is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the database instance.

* `status` - (Required, String, NonUpdatable) Specifies the switch status of the SQL limiting.  
  The valid values are as follows:
  + **ON**: Enable the SQL limiting switch.
  + **OFF**: Disable the SQL limiting switch.

* `datastore_type` - (Required, String, NonUpdatable) Specifies the database type.  
  The valid values are as follows:
  + **MySQL**: Cloud database RDS for MySQL.
  + **PostgreSQL**: Cloud database RDS for PostgreSQL.
  + **MariaDB**: Cloud database RDS for MariaDB.
  + **TaurusDB**: Cloud database TaurusDB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
