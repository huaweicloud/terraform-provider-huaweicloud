---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_custom_template"
description: ""
---

# huaweicloud_dcs_custom_template

Manages a DCS custom template resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_template_id" {}

resource "huaweicloud_dcs_custom_template" "test"{
  template_id = var.source_template_id
  name        = "template_name"
  type        = "sys"

  params {
    param_name  = "timeout"
    param_value = "100"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `template_id` - (Required, String, ForceNew) Specifies the ID of the source template.

  Changing this parameter will create a new resource.

* `source_type` - (Required, String, ForceNew) Specifies the type of the source template. Value options:
  + **sys**: system template.
  + **user**: custom template.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the template.

* `params` - (Required, List) Specifies the list of the template params.
The [params](#CustomTemplate_Param) structure is documented below.

* `description` - (Optional, String) Specifies the description of the template.

<a name="CustomTemplate_Param"></a>
The `params` block supports:

* `param_name` - (Required, String) Indicates the name of the param. You can find it through data source
  `huaweicloud_dcs_template_detail`.

* `param_value` - (Required, String) Indicates the value of the param.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `type` - Indicates the type of the template. The value can be **sys**, **user**.

* `engine` - Indicates the cache engine. Currently, only **Redis** is supported.

* `engine_version` - Indicates the cache engine version. The value can be **4.0**, **5.0.**, **6.0.**.

* `cache_mode` - Indicates the DCS instance type. The value can be **single**, **ha**, **cluster**, **proxy**,
  **ha_rw_split**.

* `product_type` - Indicates the product edition. The value can be **generic**, **enterprise**.

* `storage_type` - Indicates the storage type. The value can be **DRAM**, **SSD**.

* `params` - Indicates the list of the template params.
  The [params](#CustomTemplate_Param) structure is documented below.

* `created_at` - Indicates the time when the custom template is created.

<a name="CustomTemplate_Param"></a>
The `params` block supports:

* `param_id` - (Optional, String) Indicates the ID of the param.

* `default_value` - (Optional, String) Indicates the default value of the param.

* `value_range` - (Optional, String) Indicates the value range of the param.

* `value_type` - (Optional, String) Indicates the value type of the param.

* `description` - (Optional, String) Indicates the description of the param.

* `need_restart` - (Optional, Bool) Indicates whether the DCS instance need restart.

## Import

The dcs custom template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dcs_custom_template.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `template_id`, `source_type`, `params`. It
is generally recommended running `terraform plan` after importing a custom template. You can then decide if changes
should be applied to the custom template, or the resource definition should be updated to align with the DCS custom
template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dcs_custom_template" "test" {
  ...

  lifecycle {
    ignore_changes = [
      template_id, source_type, params,
    ]
  }
}
```
