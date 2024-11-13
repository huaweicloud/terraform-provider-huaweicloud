---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_pt_application_records"
description: |-
  Use this data source to get the list of DDS parameter template application records.
---

# huaweicloud_dds_pt_application_records

Use this data source to get the list of DDS parameter template application records.

## Example Usage

```hcl
variable "configuration_id"  {}

data "huaweicloud_dds_pt_application_records" "test" {
  configuration_id = var.configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `configuration_id` - (Required, String) Specifies the ID of the parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - Indicates the application records.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `applied_at` - Indicates the application time, in the **yyyy-mm-ddThh:mm:ssZ** format.

* `apply_result` - Indicates the application result.

* `failure_reason` - Indicates the failure reason.
