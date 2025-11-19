---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_artifact_accessories"
description: |-
  Use this data source to get the list of SWR enterprise instance artifact accessories.
---

# huaweicloud_swr_enterprise_instance_artifact_accessories

Use this data source to get the list of SWR enterprise instance artifact accessories.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}
variable "reference" {}

data "huaweicloud_swr_enterprise_instance_artifact_accessories" "test" {
  instance_id     = var.instance_id
  namespace_name  = var.namespace_name
  repository_name = var.repository_name
  reference       = var.reference
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `repository_name` - (Required, String) Specifies the repository name.

* `reference` - (Required, String) Specifies the digest.

* `type` - (Optional, String) Specifies the accessort type. Value can be **signature.cosign**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accessories` - Indicates the accessories.

  The [accessories](#accessories_struct) structure is documented below.

* `total` - Indicates the total counts of the accessories.

<a name="accessories_struct"></a>
The `accessories` block supports:

* `id` - Indicates the accessory ID.

* `type` - Indicates the accessory type.

* `size` - Indicates the accessory size.

* `digest` - Indicates the digest of the accessory.

* `artifact_id` - Indicates the artifact ID of the accessory.

* `subject_artifact_id` - Indicates the subject artifact ID of the accessory.

* `created_at` - Indicates the create time of the accessory.
