---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_log_converge_switch"
description: |-
  Using this resource to enable the LTS log receiving status within HuaweiCloud.
---

# huaweicloud_lts_log_converge_switch

Using this resource to enable the LTS log receiving status within HuaweiCloud.

~> Deleting this resource means disable the LTS log receiving status.

## Example Usage

```hcl
# Create a resource to enable the LTS log receiving status.
resource "huaweicloud_lts_log_converge_switch" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the configurations of log converge are located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
