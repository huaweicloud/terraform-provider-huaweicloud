---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_catalogs"
description: |-
  Use this data source to query application catalogs under a specified region within HuaweiCloud.
---

# huaweicloud_workspace_application_catalogs

Use this data source to query application catalogs under a specified region within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_application_catalogs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies The region where the application catalogs are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `catalogs` - The list of application catalogs.  
  The [catalogs](#workspace_application_catalogs) structure is documented below.

<a name="workspace_application_catalogs"></a>
The `catalogs` block supports:

* `id` - The ID of the application catalog.

* `zh` - The catalog description in Chinese.

* `en` - The catalog description in English.
