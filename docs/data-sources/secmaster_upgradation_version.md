---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_upgradation_version"
description: |-
  Use this data source to get the currently available version.
---

# huaweicloud_secmaster_upgradation_version

Use this data source to get the currently available version.

## Example Usage

```hcl
data "huaweicloud_secmaster_upgradation_version" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `version` - The current version.
