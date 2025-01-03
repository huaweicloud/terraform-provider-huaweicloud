---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_copy"
description: |-
  Manages a CodeArts deploy application copy resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_application_copy

Manages a CodeArts deploy application copy resource within HuaweiCloud.

## Example Usage

```hcl
variable "app_id" {}

resource "huaweicloud_codearts_deploy_application_copy" "test" {
  app_id = var.app_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `app_id` - (Required, String, ForceNew) Specifies the application ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, is same as the new application ID.
