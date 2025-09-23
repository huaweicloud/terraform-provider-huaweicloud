---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_public_scripts"
description: |-
  Use this data source to get the list of COC public scripts.
---

# huaweicloud_coc_public_scripts

Use this data source to get the list of COC public scripts.

## Example Usage

```hcl
data "huaweicloud_coc_public_scripts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `name_like` - (Optional, String) Specifies the script name, only right fuzzy search is supported.

* `name` - (Optional, String) Specifies the script name.

* `risk_level` - (Optional, String) Specifies the risk level.
  Values can be as follows:
  + **LOW**: Low risk.
  + **MEDIUM**: Medium risk.
  + **HIGH**: High risk.

* `type` - (Optional, String) Specifies the script type.
  Values can be as follows:
  + **SHELL**: Shell script.
  + **PYTHON**: Python script.
  + **BAT**: Bat script.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of public scripts.

  The [data](#data_data_struct) structure is documented below.

<a name="data_data_struct"></a>
The `data` block supports:

* `id` - Indicates the auto-increment ID of a script.

* `script_uuid` - Indicates the script UUID.

* `name` - Indicates the script name.

* `type` - Indicates the script type.

* `gmt_created` - Indicates the creation time.

* `description` - Indicates the script description.

* `properties` - Indicates the script additional properties.

  The [properties](#data_properties_struct) structure is documented below.

<a name="data_properties_struct"></a>
The `properties` block supports:

* `risk_level` - Indicates the risk level.

* `version` - Indicates the script version number.
