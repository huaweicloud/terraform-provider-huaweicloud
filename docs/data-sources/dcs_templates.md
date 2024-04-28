---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_templates"
description: ""
---

# huaweicloud_dcs_templates

Use this data source to get the list of DCS templates.

## Example Usage

```hcl
data "huaweicloud_dcs_templates" "test" {
  type = "sys"
  name = "test_template_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the template. Value options:
  + **sys**: system template.
  + **user**: custom template.

* `template_id` - (Optional, String) Specifies the ID of the template.

* `name` - (Optional, String) Specifies the name of the template.

* `engine` - (Optional, String) Specifies the cache engine. Value options: **Redis**.

* `engine_version` - (Optional, String) Specifies the cache engine version. Value options: **4.0**, **5.0**, **6.0**.

* `cache_mode` - (Optional, String) Specifies the DCS instance type. Value options:
  + **single**: single-node.
  + **ha**: master/standby.
  + **cluster**: Redis Cluster.
  + **proxy**: Proxy Cluster.
  + **ha_rw_split**: read/write splitting.

* `product_type` - (Optional, String) Specifies the product edition. Value options:
  + **generic**: standard edition.
  + **enterprise**: professional edition.

* `storage_type` - (Optional, String) Specifies the storage type. Value options: **DRAM**, **SSD**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - Indicates the list of DCS templates.
  The [templates](#Templates_Template) structure is documented below.

<a name="Templates_Template"></a>
The `templates` block supports:

* `template_id` - Indicates the ID of the template.

* `name` - Indicates the name of the template.

* `type` - Indicates the type of the template.

* `engine` - Indicates the cache engine.

* `engine_version` - Indicates the cache engine version.

* `cache_mode` - Indicates the DCS instance type.

* `product_type` - Indicates the product edition.

* `storage_type` - Indicates the storage type.

* `description` - Indicates the description of the template.
