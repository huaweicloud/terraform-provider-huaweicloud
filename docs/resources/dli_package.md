---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_package"
description: ""
---

# huaweicloud_dli_package

Manages DLI package resource within HuaweiCloud

## Example Usage

### Upload the specified python script as a resource package

```hcl
variable "group_name" {}
variable "access_domain_name" {}

resource "huaweicloud_dli_package" "queue" {
  group_name  = var.group_name
  object_path = "https://${var.access_domain_name}/dli/packages/object_file.py"
  type        = "pyFile"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to upload packages.
  If omitted, the provider-level region will be used.
  Changing this parameter will delete the current package and upload a new package.

* `group_name` - (Optional, String, ForceNew) Specifies the group name which the package belongs to.
  Changing this parameter will delete the current package and upload a new package.

* `type` - (Required, String, ForceNew) Specifies the package type.
  + **jar**: `.jar` or jar related files.
  + **pyFile**: `.py` or python related files.
  + **file**: Other user files.
  + **modelFile**: User AI model files.

  Changing this parameter will delete the current package and upload a new package.

* `object_path` - (Required, String, ForceNew) Specifies the OBS storage path where the package is located.
  For example, `https://{bucket_name}.obs.{region}.myhuaweicloud.com/dli/packages/object_file.py`.
  Changing this parameter will delete the current package and upload a new package.

* `is_async` - (Optional, Bool, ForceNew) Specifies whether to upload resource packages in asynchronous mode.
  The default value is **false**. Changing this parameter will delete the current package and upload a new package.

* `owner` - (Optional, String) Specifies the name of the package owner. The owner must be IAM user.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the package.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
  If `group_name` is specified, the ID is constructed from the `group_name` and `object_name`,
  the format is `<group_name>#<object_name>`, otherwise the resource ID which equals the `object_name`.

* `object_name` - The package name.

* `status` - Status of a package group to be uploaded.

* `created_at` - Time when a queue is created.

* `updated_at` - The last time when the package configuration update has completed.
