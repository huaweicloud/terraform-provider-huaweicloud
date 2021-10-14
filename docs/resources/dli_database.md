---
subcategory: "Data Lake Insight (DLI)"
---

# huaweicloud_dli_database

Manages DLI SQL database resource within HuaweiCloud.

## Example Usage

### Create a database

```hcl
variable "database_name" {}

resource "huaweicloud_dli_database" "test" {
  name = var.database_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DLI database resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new database resource.

* `name` - (Required, String, ForceNew) Specifies the database name. The name consists of 1 to 128 characters, starting
  with a letter or digit. Only letters, digits and underscores (_) are allowed and the name cannot be all digits.
  Changing this parameter will create a new database resource.

* `description` - (Optional, String, ForceNew) Specifies the description of a queue.
  Changing this parameter will create a new database resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  The value 0 indicates the default enterprise project. Changing this parameter will create a new database resource.

* `owner` - (Optional, String) Specifies the name of the SQL database owner.
  The owner must be IAM user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. For database resources, the ID is the database name.

## Import

DLI SQL databases can be imported by their `name`, e.g.

```
$ terraform import huaweicloud_dli_database.test terraform_test
```
