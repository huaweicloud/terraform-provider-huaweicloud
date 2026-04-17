---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_templates"
description: |-
  Use this data source to get the list of RFS templates within HuaweiCloud.
---

# huaweicloud_rfs_templates

Use this data source to get the list of RFS templates within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_rfs_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of RFS templates.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `template_id` - The unique ID of the template.

* `template_name` - The name of the template.

* `template_description` - The description of the template.

* `create_time` - The time when the template was created, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

* `update_time` - The time when the template was last updated, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

* `latest_version_description` - The description of the latest template version.

* `latest_version_id` - The ID of the latest template version.
