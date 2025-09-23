---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_dependencies"
description: |-
  Use this data source to query dependency packages within HuaweiCloud.
---

# huaweicloud_fgs_dependencies

Use this data source to query dependency packages within HuaweiCloud.

~> Between `1.64.2` and `1.72.1`, the version list of each dependency package is queried by default.
   <br>This will cause the query to take up a lot of time and may trigger flow control.
   <br>There are not recommended to use.

## Example Usage

### Obtain all public dependency packages

```hcl
data "huaweicloud_fgs_dependencies" "test" {}
```

### Obtain specific public dependency package by name

```hcl
data "huaweicloud_fgs_dependencies" "test" {
  type = "public"
  name = "obssdk-3.0.2"
}
```

### Obtain all public Python2.7 dependency packages

```hcl
data "huaweicloud_fgs_dependencies" "test" {
  type    = "public"
  runtime = "Python2.7"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region where the dependency packages are located.  
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the type of the dependency package.  
  The valid values are as follows:
  + **public**
  + **private**

* `runtime` - (Optional, String) Specifies the runtime of the dependency package.  
  The valid values are as follows:
  + **Java8**
  + **Java11**
  + **Node.js6.10**
  + **Node.js8.10**
  + **Node.js10.16**
  + **Node.js12.13**
  + **Node.js14.18**
  + **Node.js16.17**
  + **Node.js18.15**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Python3.10**
  + **Go1.x**
  + **C#(.NET Core 2.0)**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **Custom**
  + **PHP7.3**
  + **Cangjie1.0**
  + **http**
  + **Custom Image**

* `name` - (Optional, String) Specifies the name of the dependency package.

* `is_versions_query_allowed` - (Optional, Bool) Specifies whether to query the versions of each dependency package.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `packages` - All dependency packages that match the filter parameters.
  The [packages](#dependency_packages) structure is documented below.

<a name="dependency_packages"></a>
The `packages` block supports:

* `id` - The ID of the dependency package.

* `name` - The name of the dependency package.

* `owner` - The owner of the dependency package.

* `link` - The OBS bucket path where the dependency package is located (FunctionGraph serivce side).

* `etag` - The unique ID of the dependency package.

* `size` - The size of the dependency package.

* `file_name` - The file name of the stored dependency package.

* `runtime` - The runtime of the dependency package.

* `versions` - The list of the versions for the dependency package.
  The [versions](#dependency_versions) structure is documented below.

<a name="dependency_versions"></a>
The `versions` block supports:

* `id` - The ID of the dependency package version.

* `version` - The dependency package version.
