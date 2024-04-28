---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_quotas"
description: ""
---

# huaweicloud_dli_quotas

Use this data source to get a list of resource quotas for DLI service within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dli_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the resource type that used to query corresponding quota.
  Value options:
   + **CU**: Computing unit.
   + **QUEUE**: Resource queue.
   + **DATABASE**: Database.
   + **TABLE**: Table.
   + **TEMPLATE**: Template.
   + **SL_PKG_RESOURCE**: Spark job resource package.
   + **SL_SESSION**: Spark job table.
   + **JOB_CU**: Job computing unit.
   + **ELASTIC_RESOURCE_POOL**: Elastic resource pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of resource quotas.

  The [quotas](#Quotas_quotas_struct) structure is documented below.

<a name="Quotas_quotas_struct"></a>
The `quotas` block supports:

* `min` - The minimum quota of resource.

* `max` - The maximum quota of resource.

* `quota` - The current quota of resource.

* `used` - The used quota of resource.

* `type` - The resource type corresponding to quota.
