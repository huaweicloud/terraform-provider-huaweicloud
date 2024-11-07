---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_share_app"
description: |-
  Manages a CPH share app resource within HuaweiCloud.
---

# huaweicloud_cph_share_app

Manages a CPH share app resource within HuaweiCloud.

## Example Usage

```hcl
variable "package_name" {}
variable "bucket_name" {}
variable "object_path" {}
variable "server_id" {}

resource "huaweicloud_cph_share_app" "test" {
  server_id       = var.server_id
  package_name    = var.package_name
  bucket_name     = var.bucket_name
  object_path     = var.object_path
  pre_install_app = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the CPH server ID.

* `package_name` - (Required, String, NonUpdatable) Specifies the package name.
  The valid value can be **com.cph.config**, **com.cph.config.level1** or **com.cph.config.level2**.

* `bucket_name` - (Required, String, NonUpdatable) Specifies the OBS bucket name.

* `object_path` - (Required, String, NonUpdatable) Specifies the OBS object path.
  The naming format of tar file type is **<package_name>_<version_name>.tar**. For example, **com.cph.config_v1.1**.

* `pre_install_app` - (Optional, Int, NonUpdatable) Specifies whether to pre-install the application.
  The valid value can be **0** (Pre-installed) or **1** (Not pre-installed). Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `delete` - Default is 30 minutes.
