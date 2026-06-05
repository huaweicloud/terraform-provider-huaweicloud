---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_siem_shippers"
description: |-
  Use this data source to get the list of siem shippers.
---

# huaweicloud_secmaster_siem_shippers

Use this data source to get the list of siem shippers.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_siem_shippers" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `dataspace_id` - (Optional, String) Specifies the dataSpace ID.

* `pipe_id` - (Optional, String) Specifies the pipe ID.

* `shipper_name` - (Optional, String) Specifies the shipper name.

* `shipper_source_region` - (Optional, String) Specifies the shipper source region.

* `shipper_source_strategy` - (Optional, String) Specifies the shipper source strategy.

* `shipper_consumption_type` - (Optional, String) Specifies the shipper consumption type.

* `destination_shipper_type` - (Optional, String) Specifies the destination shipper type.

* `shipper_status` - (Optional, String) Specifies the shipper status.

* `create_time` - (Optional, String) Specifies the creation time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The result data list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `consumption_type` - The consumption type.

* `create_time` - The creation time, in milliseconds.

* `dataspace_id` - The dataSpace ID.

* `dataspace_name` - The dataSpace name.

* `domain_id` - The domain ID.

* `id` - The ID.

* `pipe_id` - The pipe ID.

* `pipe_name` - The pipe name.

* `project_id` - The project ID.

* `shipper_destination` - The shipper destination.

  The [shipper_destination](#shipper_destination_struct) structure is documented below.

* `shipper_id` - The shipper ID.

* `shipper_name` - The shipper name.

* `shipper_source` - The shipper source.

  The [shipper_source](#shipper_source_struct) structure is documented below.

* `status` - The status.

* `update_time` - The update time, in milliseconds.

* `version` - The version.

* `workspace_id` - The workspace ID.

<a name="shipper_destination_struct"></a>
The `shipper_destination` block supports:

* `data_param` - The data parameter, usually in JSON format.

* `data_type` - The data type.

* `dataspace` - The dataSpace ID.

* `dataspace_name` - The dataSpace name.

* `destination_info` - The destination information.

* `id` - The ID.

* `identity` - The identity.

* `pipe` - The pipe ID.

* `pipe_name` - The pipe name.

* `region` - The region.

* `type` - The type.

* `workspace` - The workspace ID.

* `workspace_name` - The workspace name.

<a name="shipper_source_struct"></a>
The `shipper_source` block supports:

* `data_type` - The data type.

* `dataspace` - The dataSpace ID.

* `dataspace_name` - The dataSpace name.

* `id` - The ID.

* `identity` - The identity.

* `pipe` - The pipe ID.

* `pipe_name` - The pipe name.

* `region` - The region.

* `type` - The type.

* `workspace` - The workspace ID.

* `workspace_name` - The workspace name.
