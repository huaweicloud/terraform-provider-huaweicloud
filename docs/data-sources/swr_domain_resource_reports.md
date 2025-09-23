---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_domain_resource_reports"
description: |-
  Use this data source to get the list of SWR domain resource reports.
---

# huaweicloud_swr_domain_resource_reports

Use this data source to get the list of SWR domain resource reports.

## Example Usage

```hcl
variable "resource_type" {}
variable "frequency" {}

data "huaweicloud_swr_domain_resource_reports" "test" {
  resource_type = var.resource_type
  frequency     = var.frequency
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The valid values are as follows:
  + **downflow**
  + **store**

* `frequency` - (Required, String) Specifies the frequency type.
  The valid value can be **daily**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `reports` - The domain resource reports.

  The [reports](#reports_struct) structure is documented below.

<a name="reports_struct"></a>
The `reports` block supports:

* `date` - The date of the domain resource report.

* `value` - The value of the domain resource report.
