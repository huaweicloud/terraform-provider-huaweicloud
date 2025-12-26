---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_addons"
description: |-
  Use this datasource to get the list of CCE Autopilot Addon Instances.
---

# huaweicloud_cce_autopilot_addons

Use this datasource to get the list of CCE Autopilot Addon Instances.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_autopilot_addons" "test" {
  cluster_id = var.cluster_id  # Required
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `addon_template_name` - (Optional, String) Specifies the name of the add-on.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The add-on instance list.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `metadata` - The metadata of the add-on instance.

  The [metadata](#items_metadata_struct) structure is documented below.

* `spec` - The detailed description of the add-on instance.

  The [spec](#items_spec_struct) structure is documented below.

* `status` - The add-on instance status.

  The [status](#items_status_struct) structure is documented below.

<a name="items_metadata_struct"></a>
The `metadata` block supports:

* `uid` - The add-on instance ID.

* `name` - The add-on instance name.

* `alias` - The add-on instance alias.

* `labels` - The add-on labels in key/value pairs.

* `annotations` - The add-on annotations in the format of key/value pairs.

* `update_timestamp` - The update time.

* `creation_timestamp` - The creation time.

<a name="items_spec_struct"></a>
The `spec` block supports:

* `addon_template_labels` - The labels of the add-on template.

* `description` - The add-on description.

* `values` - The add-on installation parameters.

* `cluster_id` - The cluster ID.

* `version` - The add-on version.

* `addon_template_name` - The add-on name.

* `addon_template_type` - The add-on type.

* `addon_template_logo` - The URL for obtaining the add-on template logo.

<a name="items_status_struct"></a>
The `status` block supports:

* `current_version` - The information about the current add-on version.

  The [current_version](#status_current_version_struct) structure is documented below.

* `is_rollbackable` - Whether the add-on version can be rolled back to the source version.

* `previous_version` - The add-on version before upgrade or rollback

* `status` - The statuses of add-on instances.

* `reason` - The cause of the add-on installation failure.

* `message` - The installation error details.

* `target_versions` - The versions to which the current add-on version can be upgraded.

<a name="status_current_version_struct"></a>
The `current_version` block supports:

* `input` - The add-on installation parameters.

* `stable` - Whether the add-on version is a stable release.

* `translate` - The translation information used by the GUI.

* `support_versions` - The cluster versions that support the add-on.

  The [support_versions](#current_version_support_versions_struct) structure is documented below.

* `creation_timestamp` - The creation time.

* `update_timestamp` - The update time.

* `version` - The add-on version.

<a name="current_version_support_versions_struct"></a>
The `support_versions` block supports:

* `cluster_type` - The cluster type that supports the add-on.

* `cluster_version` - The cluster versions that support the add-on. The value is a regular expression.

* `category` - The current support version category.
