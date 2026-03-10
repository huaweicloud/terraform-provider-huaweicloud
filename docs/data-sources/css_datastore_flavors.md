---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_datastore_flavors"
description: |-
  Use this data source to query available node specifications by engine.
---

# huaweicloud_css_datastore_flavors

Use this data source to query available node specifications by engine.

## Example Usage

```hcl
variable "datastore_id" {}

data "huaweicloud_css_datastore_flavors" "test" {
  datastore_id = var.datastore_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `datastore_id` - (Required, String) Specifies the engine type ID.

* `datastore_version_id` - (Optional, String) Specifies the engine version ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastore_id_str` - The engine type ID.

* `dbname` - The engine name.

* `versions` - The list of engine versions.
  The [versions](#versions_struct) structure is documented below.

* `model_list` - The models information.
  The [model_list](#model_list_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `id` - The engine version ID.

* `name` - The version name.

* `flavors` - The flavors information.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `cpu` - The number of vCPUs available with an instance.

* `ram` - The memory size of an instance, in GB.

* `name` - The flavor name.

* `region` - The regions where the node flavor is available.

* `typename` - The node type name.

* `diskrange` - The disk capacity of an instance, in GB.

* `cond_operation_status` - The flavor sales status.
  This parameter takes effect region-wide.
  + **normal**: The flavor is in normal commercial use.
  + **sellout**: The flavor has been sold out.

* `cond_operation_az` - The flavor status.
  This parameter takes effect AZ-wide.

* `localdisk` - Whether the node uses local disks.

* `flavor_type_cn` - The flavor categories in Chinese.

* `flavor_type_en` - The flavor categories in English.

* `edge` - Whether this is a node flavor for edge deployments.

* `str_id` - The flavor ID.

* `is_allow_https` - Whether the node type supports HTTPS access.

<a name="model_list_struct"></a>
The `model_list` block supports:

* `total_size` - The model quantity.

* `models` - The model list.
  The [models](#models_struct) structure is documented below.

<a name="models_struct"></a>
The `models` block supports:

* `id` - The model ID.

* `name` - The model name.

* `datastore_type` - The model type.

* `datastore_version` - The model version.

* `is_text_model` - Whether it is a text model.

* `model_version_id` - The model version ID.

* `desc` - The model description.

* `language` - The model language.

* `arch_type` - The model specifications.
