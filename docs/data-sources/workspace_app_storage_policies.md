---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_storage_policies"
description: |-
  Use this data source to get the list of storage permission policies for Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_storage_policies

Use this data source to get the list of storage permission policies for Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_storage_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the custom storage permission policies are located.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of the storage permission policies.  
  The [policies](#workspace_app_storage_policies) structure is documented below.

<a name="workspace_app_storage_policies"></a>
The `policies` block supports:

* `id` - The ID of the storage permission policy.

* `server_actions` - The collection of permissions that server can use to access storage.

* `client_actions` - The collection of permissions that client can use to access storage.
