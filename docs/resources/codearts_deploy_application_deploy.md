---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_deploy"
description: |-
  Manages a CodeArts deploy application deploy resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_application_deploy

Manages a CodeArts deploy application deploy resource within HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_codearts_deploy_application_deploy" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `task_id` - (Required, String, ForceNew) Specifies the deployment task ID.
  Changing this creates a new resource.

* `params` - (Optional, List, ForceNew) Specifies the parameters transferred during application deployment.
  Changing this creates a new resource.
  The [params](#block--params) structure is documented below.

* `record_id` - (Optional, String, ForceNew) Specifies the deployment record ID of an application. Specifies it to roll
  back the application to the previous deployment status.
  Changing this creates a new resource.

* `trigger_source` - (Optional, String, ForceNew) Specifies the trigger source.
  Valid values are as follows:
  + **0**: Deployment can be triggered through all requests.
  + **1**: Deployment can be triggered only through pipeline.

  Changing this creates a new resource.

<a name="block--params"></a>
The `params` block supports:

* `name` - (Optional, String, ForceNew) Specifies the parameter name transferred when deploying application.
  Changing this creates a new resource.

* `type` - (Optional, String, ForceNew) Specifies the parameter type. If a dynamic parameter is set, the type is mandatory.
  Changing this creates a new resource.

* `value` - (Optional, String, ForceNew) Specifies the parameter value transferred during application deployment.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The application deployment record can be imported using `task_id`, and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_application_deploy.test <task_id>/<id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `params`, `record_id` and `trigger_source`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the deployment record, or the resource definition should be updated to
align with the deployment record. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_deploy_application_deploy" "test" {
    ...

  lifecycle {
    ignore_changes = [
      params, record_id, trigger_source,
    ]
  }
}
```
