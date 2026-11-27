---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_associated_resource_settings"
description: |-
  Use this data source to query the TMS associated resource settings within HuaweiCloud.
---

# huaweicloud_tms_associated_resource_settings

Use this data source to query the TMS associated resource settings within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_tms_associated_resource_settings" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) Specifies the region ID used to filter settings.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `settings` - The list of associated resource settings.
  The [settings](#tms_associated_resource_settings_attr) structure is documented below.

<a name="tms_associated_resource_settings_attr"></a>
The `settings` block supports:

* `setting_name` - The name of the rule setting.

* `master_service` - The master service.

* `master_resource_type` - The master resource type.

* `associated_service` - The associated service.

* `associated_resource_type` - The associated resource type.

* `support_existing_resource` - Whether the existing resource is supported.

* `support_auto_delete` - Whether the auto delete is supported.

* `region_ids` - The list of regions supported by rule configuration.
