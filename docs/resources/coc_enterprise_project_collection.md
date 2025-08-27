---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_enterprise_project_collection"
description: |-
  Manages a COC enterprise project collection resource within HuaweiCloud.
---

# huaweicloud_coc_enterprise_project_collection

Manages a COC enterprise project collection resource within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_coc_enterprise_project_collection" "test" {
  ep_id_list = [var.enterprise_project_id]
}
```

## Argument Reference

The following arguments are supported:

* `ep_id_list` - (Required, List, NonUpdatable) Specifies the enterprise projects selected in the favorite
  configuration to form the enterprise project favorite list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
