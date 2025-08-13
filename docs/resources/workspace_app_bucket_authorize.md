---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_bucket_authorize"
description: |-
  Use this resource to create or authorize a bucket for the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_bucket_authorize

Use this resource to create or authorize a bucket for the Workspace APP within HuaweiCloud.

-> 1. If the bucket does not exist, using this resource automatically creates an OBS bucket in the specified region,
   the bucket name consists of `wks-app` and the current project ID, separated by a hyphen (-).
   e.g. `wks-app-0970dd7a1300f5672ff2c003c60ae115`.
   <br>2. This resource is a one-time operation resource used to create or authorize an OBS bucket. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_workspace_app_bucket_authorize" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the OBS bucket and the application to be
  authorized are located.  
  If omitted, the provider-level region will be used. Changing this will create new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
