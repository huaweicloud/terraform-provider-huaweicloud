---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_resources_mapping"
description: |-
  Use this data source to get the resource type mapping of EPS within HuaweiCloud.
---

# huaweicloud_enterprise_project_resources_mapping

Use this data source to get the resource type mapping of EPS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_resources_mapping" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_mapping` - The resource type mapping of EPS. This field is a key-value pair type.
