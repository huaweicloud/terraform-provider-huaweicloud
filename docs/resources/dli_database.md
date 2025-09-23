---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_database"
description: ""
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

* `name` - (Required, String, ForceNew) Specifies the database name.  
  The name consists of `1` to `128` characters, starting with a letter or digit.
  Only letters, digits and underscores (_) are allowed and the name cannot be all digits.  
  Changing this parameter will create a new database resource.

* `description` - (Optional, String, ForceNew) Specifies the description of a queue.
  Changing this parameter will create a new database resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  The value 0 indicates the default enterprise project. Changing this parameter will create a new database resource.

* `owner` - (Optional, String) Specifies the name of the SQL database owner.
  The owner must be IAM user.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the database.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID.

-> If the user has opened the EPS service, this value is a UUID value. If not, this value is the database name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

DLI SQL databases can be imported by their `name`, e.g.

```bash
$ terraform import huaweicloud_dli_database.test terraform_test
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `tags`.

It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_factory_script" "test" {
  ...

  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}
```
