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

* `support_item` - The configurations of EPS.
  The [configuration](#eps_config) structure is documented below.

<a name="eps_config"></a>
The `configuration` block supports:

* `delete_ep_support` - Whether enterprise projects can be deleted.
