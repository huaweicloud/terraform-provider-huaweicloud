---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_templates"
description: |-
  Use this data source to get the list of templates.
---

# huaweicloud_compute_templates

Use this data source to get the list of templates.

## Example Usage

```hcl
data "huaweicloud_compute_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `launch_template_id` - (Optional, List) Specifies the template IDs.

* `name` - (Optional, List) Specifies the template names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `launch_templates` - Indicates the list of templates.

  The [launch_templates](#launch_templates_struct) structure is documented below.

<a name="launch_templates_struct"></a>
The `launch_templates` block supports:

* `id` - Indicates the template ID.

* `name` - Indicates the template name.

* `description` - Indicates the template description.

* `default_version` - Indicates the default version of the template.

* `latest_version` - Indicates the latest version of the template.

* `created_at` - Indicates the time when the template was created.

* `updated_at` - Indicates the time when the template was updated.
