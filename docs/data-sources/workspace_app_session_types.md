---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_session_types"
description: |-
  Use this data source to get session type list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_session_types

Use this data source to get session type list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_session_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `session_types` - The list of session types.

  The [session_types](#workspace_app_session_types) structure is documented below.

<a name="workspace_app_session_types"></a>
The `session_types` block supports:

* `resource_spec_code` - The resource specification code.

* `session_type` - The session type.

* `resource_type` - The resource type.

* `cloud_service_type` - The code of cloud service type to which the resource belongs.
