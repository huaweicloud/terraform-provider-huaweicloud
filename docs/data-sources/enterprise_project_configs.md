---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_configs"
description: |-
  Use this data source to get the configurations of EPS within HuaweiCloud.
---

# huaweicloud_enterprise_project_configs

Use this data source to get the configurations of EPS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_configs" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `support_item_attribute` - The configurations of EPS.
  The [support_item_attribute](#support_item_attribute_struct) structure is documented below.

<a name="support_item_attribute_struct"></a>
The `support_item_attribute` block supports:

* `delete_ep_support_attribute` - Whether enterprise projects can be deleted.
