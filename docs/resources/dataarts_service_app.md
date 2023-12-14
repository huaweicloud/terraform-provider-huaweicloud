---
subcategory: "DataArts Studio"
---

# huaweicloud_dataarts_service_app

Manages an app resource of DataArts DataService within HuaweiCloud.

An app is a set of API access permissions and defines the identity of an API caller.
Each app corresponds to a unique identity credential and can be classified based on the app issuer.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_dataarts_service_app" "test" {
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

* `workspace_id` - (Required, String, ForceNew) The workspace ID.

  Changing this parameter will create a new resource.

* `dlm_type` - (Required, String, ForceNew) The type of DLM.  
  The valid values are as follows:
    - **SHARED**: Shared data service.
    - **EXCLUSIVE**: The exclusive data service.

  Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the app.  
  The name must be the account name when the `app_type` is **IAM**.

* `description` - (Optional, String) The description of the app.

* `app_type` - (Optional, String, ForceNew) The type of the app.  
  The valid values are **APP** and **IAM**.
  Defaults to **APP**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `app_key` - The key of the app.

* `app_secret` - The secret of the app.

## Import

The DataArts app can be imported using **workspace_id**, **dlm_type** and **id**, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_service_app.test <workspace_id>/<dlm_type>/<id>
```
