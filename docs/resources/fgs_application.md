---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_application"
description: ""
---

# huaweicloud_fgs_application

Manages an application within HuaweiCloud.

-> Currently, only available in `cn-north-4` and `cn-east-3` regions.

## Example Usage

### Create a simple application

```hcl
variable "application_name"
variable "application_template_id"
variable "agency_name"

resource "huaweicloud_fgs_application" "test" {
  name        = var.application_name
  template_id = var.application_template_id
  agency_name = var.agency_name
  description = "Created by terraform script"
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create an application.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the application name.  
  The name can contain a maximum of 60 characters and must start with a letter and end with a letter or digit.
  Only letters, digits, underscores (_) and hyphens (-) are allowed.  
  Changing this parameter will create a new resource.

* `template_id` - (Required, String, ForceNew) Specifies the ID of the template used by the application.  
  Changing this parameter will create a new resource.

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name used by the application.  
  Changing this parameter will create a new resource.

  -> If omitted, the service will automatically create an agency, please ensure that the tenant has IAM related
     permissions. The agency will be deleted when the application is deleted.

* `description` - (Optional, String, ForceNew) Specifies the description of the application.  
  The description can contain a maximum of `1,024` characters.  
  Changing this parameter will create a new resource.

* `params` - (Optional, String, ForceNew) Specifies the template parameters, in JSON format.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The application ID in UUID format.

* `stack_id` - The ID of the stack where the application is deployed.

* `stack_resources` - The list of the stack resources information.  
  The [stack_resources](#fgs_app_stack_resources) structure is documented below.

* `repository` - The repository information.  
  The [repository](#fgs_app_repository) structure is documented below.

* `status` - The status of the application.

* `updated_at` - The latest update time of the application.

* `apig_url` - The dependency package size in bytes.

<a name="fgs_app_stack_resources"></a>
The `stack_resources` block supports:

* `physical_resource_id` - The physical resource ID.

* `physical_resource_name` - The physical resource name.

* `logical_resource_name` - The logical resource name.

* `logical_resource_type` - The logical resource type.

* `resource_status` - The status of resource.

* `status_message` - The status information.

* `href` - The hyperlink.

* `display_name` - The cloud service name.

<a name="fgs_app_repository"></a>
The `repository` block supports:

* `https_url` - The HTTP address of the repository.

* `web_url` - The repository link.

* `status` - The repository status.

* `project_id` - The project ID of the repository.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

Application can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_fgs_application.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response.
The missing attributes include: `template_id`, `agency_name`, `params`.
It is generally recommended running `terraform plan` after importing the application.
You can then decide if changes should be applied to the application, or the resource definition should be updated to
align with the application. Also you can ignore changes as below.

```hcl
resource "huaweicloud_fgs_application" "test" {
  ...

  lifecycle {
    ignore_changes = [
      template_id, agency_name,
    ]
  }
}
```
