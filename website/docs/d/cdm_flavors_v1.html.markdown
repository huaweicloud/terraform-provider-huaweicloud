---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_flavors_v1"
sidebar_current: "docs-huaweicloud-datasource-cdm-flavors-v1"
description: |-
  Get the flavor information on a HuaweiCloud cdm service.
---

# huaweicloud\_cdm\_flavors\_v1

Use this data source to get available Huaweicloud cdm flavors.

## Example Usage

```hcl
data "huaweicloud_cdm_flavors_v1" "flavor" {
}
```

## Attributes Reference

The following attributes are exported:

* `version` -
  The version of the flavor.

* `flavors` -
  Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the cdm flavor.
* `id` - The id of the cdm flavor.
