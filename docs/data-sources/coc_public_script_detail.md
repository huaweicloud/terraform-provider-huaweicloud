---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_public_script_detail"
description: |-
  Use this data source to get the COC public script detail.
---

# huaweicloud_coc_public_script_detail

Use this data source to get the COC public script detail.

## Example Usage

```hcl
variable "script_uuid" {}

data "huaweicloud_coc_public_script_detail" "test" {
  script_uuid = var.script_uuid
}
```

## Argument Reference

The following arguments are supported:

* `script_uuid` - (Required, String) Specifies the public script UUID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - Indicates the script name.

* `description` - Indicates the script description.

* `type` - Indicates the script type.

* `content` - Indicates the script content.

* `script_params` - Indicates the script input parameters.

  The [script_params](#data_script_params_struct) structure is documented below.

* `gmt_created` - Indicates the creation time.

* `properties` - Indicates the script attachment property.

  The [properties](#data_properties_struct) structure is documented below.

<a name="data_script_params_struct"></a>
The `script_params` block supports:

* `param_name` - Indicates the parameter name.

* `param_value` - Indicates the parameter value.

* `param_description` - Indicates the parameter description.

* `param_order` - Indicates the order of parameters.

* `sensitive` - Indicates whether it is a sensitive parameter.

<a name="data_properties_struct"></a>
The `properties` block supports:

* `risk_level` - Indicates the risk level.

* `version` - Indicates the script version number.
