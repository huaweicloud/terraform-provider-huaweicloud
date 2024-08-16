---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_dataservice_app"
description: ""
---

# huaweicloud_dataarts_dataservice_app

Manages an app resource of DataArts DataService within HuaweiCloud.

An app is a set of API access permissions and defines the identity of an API caller.
Each app corresponds to a unique identity credential and can be classified based on the app issuer.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_dataarts_dataservice_app" "test" {
  workspace_id = var.workspace_id
  dlm_type     = "SHARED"
  app_type     = "APP"
  name         = "demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID.

  Changing this parameter will create a new resource.

* `dlm_type` - (Required, String, ForceNew) Specifies the type of DLM engine.  
  The valid values are as follows:
  + **SHARED**: Shared data service.
  + **EXCLUSIVE**: The exclusive data service.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the application.  
  The name must be the **account** name when the `app_type` is **IAM**.

* `description` - (Optional, String) Specifies the description of the application.

* `app_type` - (Optional, String, ForceNew) Specifies the type of the application.  
  The valid values are as follows:
  + **APP**: access through app authentication.
  + **IAM**: IAM authentication is used, which means access using a token.

  Defaults to **APP**. Changing this parameter will create a new resource.

  -> The IAM app can only have one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `app_key` - The key of the app.

* `app_secret` - The secret of the app.

## Import

The DataArts DataService app can be imported using `workspace_id`, `dlm_type` and `id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_dataarts_dataservice_app.test <workspace_id>/<dlm_type>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `app_type`.
It is generally recommended running `terraform plan` after importing an application.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the application. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_dataservice_app" "test" {
  ...

  lifecycle {
    ignore_changes = [
      app_type,
    ]
  }
}
```
