---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_original_policy"
description: |-
  Use this data source to get the original policy of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_original_policy

Use this data source to get the original policy of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_original_policy" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the original policy is located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policy` - The original policy configuration, in JSON format.
