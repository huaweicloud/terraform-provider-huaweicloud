---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_persistent_volume_claims"
description: |-
  Use this data source to get the list of CCI persistent volume claim resource within HuaweiCloud.
---

# huaweicloud_cciv2_persistent_volume_claims

Use this data source to get the list of CCI persistent volume claim resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_persistent_volume_claims" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pvcs` - The persistent volume claims.
  The [pvcs](#pvcs) structure is documented below.

<a name="pvcs"></a>
The `pvcs` block supports:

* `access_modes` - The access modes of the persistent volume claim.

* `annotations` - The annotations of the persistent volume claim.

* `creation_timestamp` - The creation timestamp of the persistent volume claim.

* `finalizers` - The finalizers of the persistent volume claim.

* `labels` - The labels of the persistent volume claim.

* `name` - The name of the persistent volume claim in the namespace.

* `namespace` - The namespace.

* `resource_version` - The resource version of the persistent volume claim.

* `resources` - The access_modes of the persistent volume claim.
  The [resources](#pvcs_resources) structure is documented below.

* `selector` - The selector of the persistent volume claim.
  The [selector](#pvcs_selector) structure is documented below.

* `status` - The status of the persistent volume claim.

* `storage_class_name` - The storage class name of the persistent volume claim.

* `uid` - The uid of the persistent volume claim.

* `valume_name` - The valume name of the persistent volume claim.

* `volume_mode` - The volume mode of the persistent volume claim.

<a name="pvcs_resources"></a>
The `resources` block supports:

* `limits` - The limits expressions of the resources.

* `requests` - The requests labels of the resources.

<a name="pvcs_selector"></a>
The `selector` block supports:

* `match_expressions` - The match expressions of the selector.
  The [match_expressions](#pvcs_selector_match_expressions) structure is documented below.

* `match_labels` - The match labels of the selector.

<a name="pvcs_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match expressions.

* `operator` - The operator of the match expressions.

* `values` - The values of the match expressions.
