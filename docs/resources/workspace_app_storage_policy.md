---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_storage_policy"
description: |-
  Manages the custom storage permission policy of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_storage_policy

Manages the custom storage permission policy of Workspace APP within HuaweiCloud.

-> **NOTE:** Deleting this resource will not initialize (restore) the permission policy configuration and just only
   remove the tfstate record for this resource.

## Example Usage

### Read and write for server permission and allow upload and download permissions for client

```hcl
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject", "PutObject", "DeleteObject"]
  client_actions = ["PutObject", "DeleteObject"]
}
```

### Read-Only for server permission and deny upload and download permissions for client

```hcl
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the custom storage permission policy is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `server_actions` - (Required, List) Specifies the collection of permissions that server can use to access storage.  
  The valid configuration combinations are as follows:
  + **GetObject**: Read-Only.
  + **PutObject + GetObject + DeleteObject**: Read and Write.

* `client_actions` - (Optional, List) Specifies the collection of permissions that client can use to access storage.  
  The valid values are as follows:
  + **GetObject**: Download permission only.
  + **PutObject + DeleteObject**: Upload permission only.
  + **GetObject + PutObject + DeleteObject**: Both download and upload permissions are allowed.

  If omitted, both download and upload permissions are denied.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Custom storage permission policy can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_storage_policy.test <id>
```

'NA' or other characters can be used to instead of the `id`.

```bash
$ terraform import huaweicloud_workspace_app_storage_policy.test NA
```
