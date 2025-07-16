---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_group_restrict"
description: |-
  Use this data source to get the restrict information of Workspace APP server group within HuaweiCloud.
---

# huaweicloud_workspace_app_server_group_restrict

Use this data source to get the restrict information of Workspace APP server group within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_server_group_restrict" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `max_session` - The maximum number of connection sessions per server.

* `max_group_count` - The maximum number of server groups that can be created by the tenant.
