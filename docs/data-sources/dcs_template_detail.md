---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_template_detail"
description: ""
---

# huaweicloud_dcs_template_detail

Use this data source to get the detail of DCS template.

## Example Usage

```hcl
variable "template_id" {}
data "huaweicloud_dcs_template_detail" "test" {
  template_id = var.template_id
  type        = "sys"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the template. Value options:
  + **sys**: system template.
  + **user**: custom template.

* `template_id` - (Required, String) Specifies the ID of the template.

* `params` - (Optional, List) Specifies the list of the template params.
The [params](#TemplateDetail_Param) structure is documented below.

<a name="TemplateDetail_Param"></a>
The `params` block supports:

* `param_name` - (Optional, String) Specifies the name of the param.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - Indicates the name of the template.

* `type` - Indicates the type of the template. The value can be **sys**, **user**.

* `engine` - Indicates the cache engine. Currently, only **Redis** is supported.

* `engine_version` - Indicates the cache engine version. The value can be **4.0**, **5.0**, **6.0**.

* `cache_mode` - Indicates the DCS instance type. The value can be **single**, **ha**, **cluster**, **proxy**,
  **ha_rw_split**.

* `product_type` - Indicates the product edition. The value can be **generic**, **enterprise**.

* `storage_type` - Indicates the storage type. The value can be **DRAM**, **SSD**.

* `description` - Indicates the description of the template.

* `params` - Indicates the list of the template params.
  The [params](#TemplateDetail_Param) structure is documented below.

<a name="TemplateDetail_Param"></a>
The `params` block supports:

* `param_id` - Indicates the ID of the param.

* `default_value` - Indicates the default of the param.

* `value_range` - Indicates the value range of the param.

* `value_type` - Indicates the value type of the param.

* `description` - Indicates the description of the param.

* `need_restart` - Indicates whether the DCS instance need restart.
