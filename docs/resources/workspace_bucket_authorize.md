---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_bucket_authorize"
description: |-
  Use this resource to create and authorize a default bucket within HuaweiCloud.
---

# huaweicloud_workspace_bucket_authorize

Use this resource to create and authorize a default bucket within HuaweiCloud.

-> This resource is only a one-time action resource for creating and authorizing a default bucket. Deleting this
   resource will not delete the corresponding default bucket from the OBS service, but will only remove the resource
   information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_workspace_bucket_authorize" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the bucket to be authorized is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `bucket_name` - The name of the bucket that was created and authorized.
