---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_supported_plugins"
description: |-
  Use this data source to query the supported plugins for GaussDB instances within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_supported_plugins

Use this data source to query the supported plugins for GaussDB instances within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_instance_supported_plugins" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - The list of supported plugin names.
