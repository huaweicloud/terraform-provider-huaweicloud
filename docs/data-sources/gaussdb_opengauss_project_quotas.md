---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_project_quotas"
description: |-
  Use this data source to get the project quotas of a specified tenant.
---

# huaweicloud_gaussdb_opengauss_project_quotas

Use this data source to get the project quotas of a specified tenant.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_project_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the resource type used to filter quotas. Value options: **instance**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the instance quota of a tenant.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `resources` - Indicates the resource objects.

  The [resources](#quotas_resources_struct) structure is documented below.

<a name="quotas_resources_struct"></a>
The `resources` block supports:

* `type` - Indicates the quota of a specified type.

* `used` - Indicates the number of created resources.

* `quota` - Indicates the maximum resource quota.
