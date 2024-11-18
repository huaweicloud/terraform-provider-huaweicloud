---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_quotas"
description: |-
  Use this data source to get the list of DDS quotas.
---

# huaweicloud_dds_quotas

Use this data source to get the list of DDS quotas.

## Example Usage

```hcl
data "huaweicloud_dds_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - Indicates the quotas information.
  The [quotas](#attrblock--quotas) structure is documented below.

<a name="attrblock--quotas"></a>
The `quotas` block supports:

* `mode` - Indicates the instance type.

* `quota` - Indicates the existing quota.

* `type` - Indicates the quota resource type.

* `used` - Indicates the number of the used instances.
