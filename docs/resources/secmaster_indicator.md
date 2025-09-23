---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_indicator"
description: |-
  Manages a SecMaster indicator resource within HuaweiCloud.
---

# huaweicloud_secmaster_indicator

Manages a SecMaster indicator resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "type_id" {}

resource "huaweicloud_secmaster_indicator" "test" {
  workspace_id = var.workspace_id
  name         = "demo"
  
  type {
    category       = "Domain"
    indicator_type = "Domain"
    id             = var.type_id
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Open"
  confidence            = 80
  first_occurrence_time = "2023-10-24T17:23:55.000+08:00"
  last_occurrence_time  = "2023-10-25T11:15:30.000+08:00"
  threat_degree         = "Black"
  granularity           = "1"
  value                 = "test.terraform.com"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the indicator belongs.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the indicator name.

* `type` - (Required, List) Specifies the indicator type.
  The [type](#Indicator_IndicatorType) structure is documented below.

* `threat_degree` - (Required, String) Specifies the threat degree.
  The value can be: **Black**, **White** and **Gray**.

* `data_source` - (Required, List, ForceNew) Specifies the data source of the indicator.
  The [data_source](#Indicator_DataSource) structure is documented below.

  Changing this parameter will create a new resource.

* `status` - (Required, String) Specifies the indicator status.
  The value can be: **Open**, **Closed** and **Revoked**.

* `confidence` - (Required, Int) Specifies the confidence of the indicator.
  The value ranges from `80` to `100`.

* `first_occurrence_time` - (Required, String) Specifies the first occurrence time of the indicator.
  For example: 2023-04-18T13:00:00.000+08:00

* `last_occurrence_time` - (Required, String) Specifies the last occurrence time of the indicator.
  For example: 2023-04-18T13:00:00.000+08:00

* `granularity` - (Required, Int) Specifies the granularity of the indicator.
  The value can be:
  + **1**: First time observed;
  + **2**: In-house data;
  + **3**: To be purchased;
  + **4**: Queried from external networks;

* `value` - (Required, String) Specifies the value of the indicator.

* `labels` - (Optional, String) Specifies the labels of the indicator in comma-separated string.

* `invalid` - (Optional, Bool) Specifies whether the indicator is invalid.

<a name="Indicator_IndicatorType"></a>
The `type` block supports:

* `category` - (Required, String) Specifies the category.

* `indicator_type` - (Required, String) Specifies the indicator type.

* `id` - (Required, String) Specifies the indicator type ID.

<a name="Indicator_DataSource"></a>
The `data_source` block supports:

* `source_type` - (Required, Int, ForceNew) Specifies the data source type.
  Changing this parameter will create a new resource.

* `product_name` - (Required, String, ForceNew) Specifies the product name.
  Changing this parameter will create a new resource.

* `product_feature` - (Required, String, ForceNew) Specifies the product feature.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The created time.

* `updated_at` - The updated time.

## Import

The indicator can be imported using the workspace ID and the indicator ID, e.g.

```bash
$ terraform import huaweicloud_secmaster_indicator.test <workspace_id>/<id>
```
