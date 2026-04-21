---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_provider_versions"
description: |-
  Use this datasource to get the list of private provider versions.
---

# huaweicloud_rfs_private_provider_versions

Use this datasource to get the list of private provider versions.

## Example Usage

```hcl
variable "provider_name" {}

data "huaweicloud_rfs_private_provider_versions" "test" {
  provider_name = var.provider_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `provider_name` - (Required, String) Specifies the name of the private provider.

* `provider_id` - (Optional, String) Specifies the ID of the private provider.

* `sort_key` - (Optional, String) Specifies the sort field, only supports **create_time**.

* `sort_dir` - (Optional, String) Specifies the ascending or descending order.
  Valid values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `versions` - The list of private provider versions. By default, these versions are sorted in descending order of the
  creation time.

  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `provider_id` - The unique ID of the private provider.

* `provider_name` - The name of the private provider.

* `provider_version` - The version number of the private provider.

* `version_description` - The description of the private provider version.

* `function_graph_urn` - The URN of the FunctionGraph function.

* `create_time` - The creation time of a private provider version. It is represented in UTC format
  (YYYY-MM-DDTHH:mm:ss.SSSZ), such as **1970-01-01T00:00:00.000Z**.
