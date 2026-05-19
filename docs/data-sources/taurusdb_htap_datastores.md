---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_datastores"
description: |-
  Use this data source to query the list of TaurusDB HTAP engine resources within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_datastores

Use this data source to query the list of TaurusDB HTAP engine resources within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_taurusdb_htap_datastores" "test" {
  engine_name = "star-rocks"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP engine resources.
  If omitted, the provider-level region will be used.

* `engine_name` - (Required, String) Specifies the HTAP engine type. Value options: **star-rocks**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastores` - The list of HTAP engine database versions.
  The [datastores](#taurusdb_htap_datastores_attr) structure is documented below.

<a name="taurusdb_htap_datastores_attr"></a>
The `datastores` block supports:

* `id` - The ID of the database version, which is unique.

* `name` - The database version number. Only the two-number major version is returned.

* `kernel_version` - The database kernel version. A complete four-number version is returned.
