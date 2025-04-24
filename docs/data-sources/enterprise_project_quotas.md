---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_quotas"
description: |-
  Use this data source to get the list of quotas of EPS resources within HuaweiCloud.
---

# huaweicloud_enterprise_project_quotas

Use this data source to get the list of quotas of EPS resources within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource quotas.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of the resource quotas.
  The [quotas](#eps_quotas) structure is documented below.

<a name="eps_quotas"></a>
The `quotas` block supports:

* `quota` - The total number of the resource quota.

* `type` - The resource type corresponding to quota.

* `used` - The used quota number.
