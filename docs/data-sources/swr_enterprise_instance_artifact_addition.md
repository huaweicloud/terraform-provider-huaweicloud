---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_artifact_addition"
description: |-
  Use this data source to get the list of SWR enterprise instance artifact addition infos.
---

# huaweicloud_swr_enterprise_instance_artifact_addition

Use this data source to get the list of SWR enterprise instance artifact addition infos.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}
variable "reference" {}
variable "addition" {}

data "huaweicloud_swr_enterprise_instance_artifact_addition" "test" {
  instance_id     = var.instance_id
  namespace_name  = var.namespace_name
  repository_name = var.repository_name
  reference       = var.reference
  addition        = var.addition
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `repository_name` - (Required, String) Specifies the repository name.

* `reference` - (Required, String) Specifies the artifact digest.

* `addition` - (Required, String) Specifies the addition info type.
  Value can be **build_history**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `build_histories` - Indicate the build histories.

  The [build_histories](#build_histories_struct) structure is documented below.

* `total` - Indicate the tital build histories.

<a name="build_histories_struct"></a>
The `build_histories` block supports:

* `media_type` - Indicate the media type.

* `size` - Indicate the size.

* `digest` - Indicate the digest.

* `created` - Indicate the create time.

* `created_by` - Indicate the create command.

* `empty_layer` - Indicate whether it is the empty layer.
