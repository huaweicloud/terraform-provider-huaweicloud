---
subcategory: "FunctionGraph"
---

# huaweicloud_fgs_applications

Use this data source to get the list of applications within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_fgs_applications" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the applications are located.  
  If omitted, the provider-level region will be used.

* `application_id` - (Optional, String) Specifies the application ID used to query specified application.

* `name` - (Optional, String) Specifies the application name used to query specified application.

* `status` - (Optional, String) Specifies the status of the application to be queried.  
  The valid values are as follows:
  + **success**: The application created successfully.
  + **repoFail**: The application repository creation failed.

* `description` - (Optional, String) Specifies the description of the application to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - All applications that match the filter parameters.
  The [applications](#applications_struct) structure is documented below.

<a name="applications_struct"></a>
The `applications` block supports:

* `id` - The ID of application.

* `name` - The name of application.

* `status` -  The status of application.

* `description` - The description of application.

* `updated_at` - The latest update time of the application, in RFC3339 format.
