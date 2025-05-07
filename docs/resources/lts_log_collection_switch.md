---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_log_collection_switch"
description: |-
  Use this resource to enable or disable log collection beyond free quota within HuaweiCloud.
---

# huaweicloud_lts_log_collection_switch

Use this resource to enable or disable log collection beyond free quota within HuaweiCloud.

-> This resource is only a one-time action resource for enabling or disabling log collection switch. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
resource "huaweicloud_lts_log_collection_switch" "test" {
  action = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the operation type of the log collection switch.  
  The valid values are as follows:
  + **enable**
  + **disable**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
