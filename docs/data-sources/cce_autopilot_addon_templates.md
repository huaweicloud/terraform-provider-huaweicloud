---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_addon_templates"
description: |-
  Use this data source to get the list of CCE Autopilot add-on templates.
---

# huaweicloud_cce_autopilot_addon_templates

Use this data source to get the list of CCE Autopilot add-on templates.

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_addon_templates" "test" {
  addon_template_name = "log-agent"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `addon_template_name` - (Optional, String) Specifies the name of the add-on.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of add-on templates.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `alias` - The add-on alias.

* `annotations` - The add-on annotations in the format of key-value pairs.

* `id` - The ID of the add-on template.

* `name` - The name of the add-on.

* `description` - The description of the add-on.

* `versions` - The versions of the add-on.

  The [versions](#spec_versions_struct) structure is documented below.

* `type` - The type of the add-on template.

* `require` - Whether the add-on is required.

* `labels` - The labels of the add-on.

<a name="spec_versions_struct"></a>
The `versions` block supports:

* `input` - The install parameters of the add-on.

* `stable` - Whether the version is stable.

* `support_versions` - The list of supported cluster versions.

  The [support_versions](#versions_support_versions_struct) structure is documented below.

* `version` - The add-on version.

<a name="versions_support_versions_struct"></a>
The `support_versions` block supports:

* `cluster_type` - The supported cluster type.

* `cluster_version` - The supported cluster version.
