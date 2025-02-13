---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_dependency_version"
description: |-
  Manages a custom dependency version within HuaweiCloud.
---

# huaweicloud_fgs_dependency_version

Manages a custom dependency version within HuaweiCloud.

-> We recommend using this resource to replace the `huaweicloud_fgs_dependency` resource for managing dependency
packages. You can migrate smoothly because the parameter behavior of the two resources is consistent.

## Example Usage

### Create a custom dependency version using an OBS bucket path where the ZIP file is located

```hcl
variable "dependency_name" {}
variable "custom_dependency_location" {}

resource "huaweicloud_fgs_dependency_version" "test" {
  name    = var.dependency_name
  runtime = "Python3.6"
  link    = var.custom_dependency_location
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region where the custom dependency version is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `runtime` - (Required, String, ForceNew) Specifies the runtime of the custom dependency version.
  The valid values are as follows:
  + **Java8**
  + **Java11**
  + **Node.js6.10**
  + **Node.js8.10**
  + **Node.js10.16**
  + **Node.js12.13**
  + **Node.js14.18**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Go1.8**
  + **Go1.x**
  + **C#(.NET Core 2.0)**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **Custom**
  + **PHP 7.3**
  + **http**

  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the custom dependency package to which the version
  belongs.  
  The name can contain a maximum of `96` characters and must start with a letter and end with a letter or digit.  
  Only letters, digits, underscores (_), periods (.), and hyphens (-) are allowed.  
  Changing this will create a new resource.

* `link` - (Required, String, ForceNew) Specifies the OBS bucket path where the dependency package is located.  
  The OBS object URL must be in ZIP format, such as
  `https://obs-terraform.obs.cn-north-4.myhuaweicloud.com/huaweicloudsdkcore.zip`.  
  Changing this will create a new resource.

  -> A link can only be used to create at most one dependency package.

* `description` - (Optional, String, ForceNew) Specifies the description of the custom dependency version.  
  The description can contain a maximum of `512` characters.  
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consists of dependency ID and version number, separated by a slash (/).  
  The format is `<name>/<version>`.

* `version` - The dependency package version.

* `version_id` - The ID of the dependency package version.

* `owner` - The dependency owner, **public** indicates a public dependency.

* `etag` - The unique ID of the dependency.

* `size` - The dependency size, in bytes.

* `dependency_id` - The ID of the dependency package corresponding to the version.

## Import

Dependency version can be imported using the resource `id`, e.g.

```bash
$ terraform import huaweicloud_fgs_dependency_version.test <id>
```

Or using related dependency package `name` and the `version` number, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_fgs_dependency_version.test <name>/<version>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `link`.
It is generally recommended running `terraform plan` after importing a dependency package.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the dependency package. Also you can ignore changes as below.

```hcl
resource "huaweicloud_fgs_dependency_version" "test" {
  ...

  lifecycle {
    ignore_changes = [
      link,
    ]
  }
}
```
