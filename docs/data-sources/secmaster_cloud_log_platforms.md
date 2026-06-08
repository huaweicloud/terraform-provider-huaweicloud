---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_cloud_log_platforms"
description: |-
  Use this data source to get the list of SecMaster cloud log platforms.
---

# huaweicloud_secmaster_cloud_log_platforms

Use this data source to get the list of SecMaster cloud log platforms.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_cloud_log_platforms" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `platforms` - The cloud log platform list.

  The [platforms](#platforms_struct) structure is documented below.

<a name="platforms_struct"></a>
The `platforms` block supports:

* `tenant_managed_domain_id` - The tenant managed domain ID.

* `platform_managed_domain_id` - The platform managed domain ID.

* `dw_region` - The data warehouse region.

* `create_time` - The creation time.

* `update_time` - The update time.

* `publish_status` - The publish status.

* `white_list` - Whether the whitelist is enabled.
