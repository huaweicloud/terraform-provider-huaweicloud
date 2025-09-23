---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_quotas"
description: |-
  Use this data source to get the resource quotas in a project.
---

# huaweicloud_rds_quotas

Use this data source to get the resource quotas in a project.

## Example Usage

```hcl
data "huaweicloud_rds_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the objects in the quota list.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - Indicates the resource list objects.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `quota` - Indicates the project resource quota.

* `used` - Indicates the number of used resources.

* `type` - Indicates the project resource type. The value can be **instance**.
