---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_attacked_top"
description: |-
  Use this data source to get the top attacked assets information within HuaweiCloud.
---

# huaweicloud_dsc_attacked_top

Use this data source to get the top attacked assets information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_attacked_top" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attacked_api_num` - The number of attacked APIs.

* `attacked_api_top` - The top attacked API information.

  The [attacked_api_top](#attacked_api_top_struct) structure is documented below.

* `attacked_asset_num` - The number of attacked assets.

* `attacked_asset_top` - The top attacked asset information.

  The [attacked_asset_top](#attacked_asset_top_struct) structure is documented below.

<a name="attacked_api_top_struct"></a>
The `attacked_api_top` block supports:

* `api_name` - The API name.

* `application_type` - The application type.

* `attacked_num` - The number of attacks.

* `instance_id` - The instance ID.

<a name="attacked_asset_top_struct"></a>
The `attacked_asset_top` block supports:

* `attacked_num` - The number of attacks.

* `db_ip` - The database IP.

* `db_name` - The database name.

* `db_type` - The database type.

* `instance_id` - The instance ID.

* `instance_name` - The instance name.
