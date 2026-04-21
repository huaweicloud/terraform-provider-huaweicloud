---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_template_versions"
description: |-
  Use this data source to get all versions of an RFS template.
---

# huaweicloud_rfs_template_versions

Use this data source to get all versions of an RFS template.

## Example Usage

```hcl
variable "template_name" {}

data "huaweicloud_rfs_template_versions" "test" {
  template_name = var.template_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `template_name` - (Required, String) Specifies the template name.

* `template_id` - (Optional, String) Specifies the template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - The list of template versions.

  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `version_id` - The version ID, such as **V1**, **V2**.

* `template_id` - The template ID.

* `template_name` - The template name.

* `version_description` - The version description.

* `create_time` - The time when the version was created, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.
