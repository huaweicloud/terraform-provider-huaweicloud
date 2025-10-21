---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_apps"
description: |-
  Use this data source to get the list of CPCS applications.
---

# huaweicloud_cpcs_apps

Use this data source to get the list of CPCS applications.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "app_name" {}

data "huaweicloud_cpcs_apps" "test" {
  app_name = var.app_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `app_name` - (Optional, String) Specifies the name of the application.

* `vpc_name` - (Optional, String) Specifies the VPC name to which the application belongs.

* `sort_key` - (Optional, String) Specifies the sort attribute.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The default value is **DESC**. Valid values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apps` - Indicates the applications list.
  The [apps](#CPCS_apps) structure is documented below.

<a name="CPCS_apps"></a>
The `apps` block supports:

* `app_id` - The application ID.

* `app_name` - The application name.

* `vpc_id` - The VPC ID to which the application belongs.

* `vpc_name` - The VPC name to which the application belongs.

* `subnet_id` - The subnet ID to which the application belongs.

* `subnet_name` - The subnet name to which the application belongs.

* `domain_id` - The account ID.

* `description` - The application description.

* `create_time` - The creation time of the application, in UNIX timestamp format.
