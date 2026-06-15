---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_siem_shipper_authorizations"
description: |-
  Use this data source to get the list of siem shipper authorizations.
---

# huaweicloud_secmaster_siem_shipper_authorizations

Use this data source to get the list of siem shipper authorizations.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_siem_shipper_authorizations" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `source_region` - (Optional, String) Specifies the source region.

* `destination_shipper_type` - (Optional, String) Specifies the destination shipper type.

* `shipper_status` - (Optional, String) Specifies the shipper status.

* `auth_status` - (Optional, String) Specifies the authorization status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The result data list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `authorize_status` - The authorization status.

* `data_source_query` - The data source query string.

* `data_type` - The data type.

* `dataspace` - The dataSpace ID.

* `id` - The ID.

* `pipe` - The pipe ID.

* `region` - The region.

* `request_time` - The request time, in milliseconds.

* `handler_time` - The authorization time, in milliseconds.

* `run_status` - The run status.

* `shipper_id` - The shipper ID.

* `shipper_name` - The shipper name.

* `workspace` - The workspace ID.
