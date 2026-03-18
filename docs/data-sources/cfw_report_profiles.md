---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_report_profiles"
description: |-
  Use this data source to get the list of report profiles within HuaweiCloud.
---

# huaweicloud_cfw_report_profiles

Use this data source to get the list of report profiles within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_report_profiles" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `category` - (Optional, String) Specifies the report type.  
  The valid values are as follows:
  + **daily**: Daily report.
  + **weekly**: Weekly report.
  + **custom**: Custom report.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The list of report profile records.

  The [data](#data_struct) structure is documented below.

* `total` - The total number of report profiles.

<a name="data_struct"></a>
The `data` block supports:

* `profile_id` - The template ID.

* `name` - The template name.

* `category` - The template type.

* `status` - The enable status.

* `report_id` - The ID of the latest report.

* `last_time` - The generation time of the latest report, in milliseconds.
