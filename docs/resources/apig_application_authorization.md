---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_authorization"
description: |-
  Using this resource to authorize APIs for application, allowing it to access the published APIs within HuaweiCloud.
---

# huaweicloud_apig_application_authorization

Using this resource to authorize APIs for application, allowing it to access the published APIs within HuaweiCloud.

~> Before binding the API(s), please make sure all APIs have been published, otherwise you will receive the following
   warning message: `Warning: Resource not found`.

-> If this resource was imported and no changes were deployed before deletion (a change must be triggered to apply the
   `api_ids` configured in the script), terraform will delete all bound APIs for current configured application in
   specified publish environment. Otherwise, terraform will only delete the bound API(s) managed by the last change.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_id" {}
variable "published_env_id" {}
variable "published_api_ids" {
  type = list(string)
}

resource "huaweicloud_apig_api_publishment" "test" {
  count = length(var.published_api_ids)

  instance_id = var.instance_id
  api_id      = var.published_api_ids[count.index]
  env_id      = var.published_env_id
}

resource "huaweicloud_apig_application_authorization" "test" {
  depends_on = [huaweicloud_apig_api_publishment.test]

  instance_id    = var.instance_id
  application_id = var.application_id
  env_id         = var.published_env_id
  api_ids        = var.published_api_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application and APIs are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  and APIs belong.  
  Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of the application authorized to access the APIs.
  Changing this will create a new resource.

* `env_id` - (Required, String, ForceNew) Specifies the environment ID where the APIs were published.
  Changing this will create a new resource.

* `api_ids` - (Required, List) Specifies the authorized API IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is `<instance_id>/<env_id>/<application_id>`.

## Import

Authorize relationships of application can be imported using `id`, separated by the slashes, e.g.

```bash
$ terraform import huaweicloud_apig_application_authorization.test <instance_id>/<env_id>/<application_id>
```
