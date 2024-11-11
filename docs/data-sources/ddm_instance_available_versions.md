---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_available_versions"
description: |-
  Use this data source to get the available versions of the DDM instance.
---

# huaweicloud_ddm_instance_available_versions

Use this data source to get the available versions of the DDM instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_ddm_instance_available_versions" test {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DDM instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `current_favored_version` - Indicates the preferred version of the current series.

* `current_version` - Indicates the current version.

* `latest_version` - Indicates the latest version.

* `previous_version` - Indicates the previous version of the current instance.

* `versions` - Indicates the available version.
