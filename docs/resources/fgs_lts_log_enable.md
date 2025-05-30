---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_lts_log_enable"
description: |-
  Use this resource to enable LTS logs for FunctionGraph within HuaweiCloud.
---

# huaweicloud_fgs_lts_log_enable

Use this resource to enable LTS logs for FunctionGraph within HuaweiCloud.

-> This resource is only a one-time action resource for enabling LTS logs for FunctionGraph. Deleting this resource will
   not disable the LTS logs, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_fgs_lts_log_enable" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the LTS log function is to be enabled.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
