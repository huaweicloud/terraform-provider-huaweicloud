---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_search_path_switch"
description: |-
  Use this resource to enable or disable the search path switch within HuaweiCloud.
---

# huaweicloud_das_search_path_switch

Use this resource to enable or disable the search path switch within HuaweiCloud.

-> This resource is a one-time action resource for switching the search path switch. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.

-> This resource only supports to switch the search path switch of **PostgreSql** instances.

## Example Usage

### Enable search path switch

```hcl
variable "connection_id" {}

resource "huaweicloud_das_search_path_switch" "test" {
  connection_id = var.connection_id
  switch_on     = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the search path switch is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `connection_id` - (Required, String, NonUpdatable) Specifies the ID of the database connection (DB user ID).

* `switch_on` - (Required, Bool) Whether to enable the search path switch.  
  The valid values are as follows:
  + **true**: Enable the search path switch.
  + **false**: Disable the search path switch.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
