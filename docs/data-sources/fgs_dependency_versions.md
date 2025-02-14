---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_dependency_versions"
description: |-
  Use this data source to get the list of dependency package versions within HuaweiCloud.
---

# huaweicloud_fgs_dependency_versions

Use this data source to get the list of dependency package versions within HuaweiCloud.

## Example Usage

### Query all versions under a specified dependency package

```hcl
variable "dependency_id" {}

data "huaweicloud_fgs_dependency_versions" "test" {
  dependency_id = var.dependency_id
}
```

### Query a specified dependency package version

```hcl
variable "dependency_id" {}
variable "version_name" {}

data "huaweicloud_fgs_dependency_versions" "test" {
  dependency_id = var.dependency_id
  version       = var.version_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the dependency package and the versions are located.  
  If omitted, the provider-level region will be used.

* `dependency_id` - (Required, String) Specifies the ID of the dependency package to which the versions belong.

* `version_id` - (Optional, String) Specifies the ID of the dependency package version.

* `version` - (Optional, Int) Specifies the version of the dependency package.

* `runtime` - (Optional, String) Specifies the runtime of the dependency package version.  
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - All dependency package versions that match the filter parameters.  
  The [versions](#dependency_versions) structure is documented below.

<a name="dependency_versions"></a>
The `versions` block supports:

* `id` - The ID of the dependency package version.

* `version` - The dependency package version.

* `dependency_id` - The ID of the dependency package corresponding to the version.

* `dependency_name` - The name of the dependency package corresponding to the version.

* `runtime` - The runtime of the dependency package version.

* `link` - The OBS bucket path where the dependency package version is located.

* `size` - The size of the ZIP file used by the dependency package version, in bytes.

* `etag` - The unique ID of the dependency.

* `owner` - The dependency owner, `public` indicates a public dependency.

* `description` - The description of the dependency package version.
