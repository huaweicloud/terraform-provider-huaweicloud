---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_persistent_volumes"
description: |-
  Use this data source to get the list of CCI persistent volumes within HuaweiCloud.
---

# huaweicloud_cciv2_persistent_volumes

Use this data source to get the list of CCI persistent volumes within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cciv2_persistent_volumes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - The name.

* `persistent_volumes` - The persistent volumes.
  The [persistent_volumes](#persistent_volumes) structure is documented below.

<a name="persistent_volumes"></a>
The `persistent_volumes` block supports:

* `annotations` - The annotations.

* `claim_ref` - The claim reference.
  The [claim_ref](#persistent_volumes_claim_ref) structure is documented below.

* `creation_timestamp` - The creation time.

* `csi` - The CSI.
  The [csi](#persistent_volumes_csi) structure is documented below.

* `finalizers` - The finalizers.

* `labels` - The labels.

* `mount_options` - The mount options.

* `node_affinity` - The node affinity.
  The [node_affinity](#persistent_volumes_node_affinity) structure is documented below.

* `reclaim_policy` - The reclaim policy.

* `resource_version` - The resource version.

* `status` - The status.
  The [status](#persistent_volumes_status) structure is documented below.

* `storage_class_name` - The storage class name.

* `uid` - The uid.

* `volume_mode` - The volume mode.

<a name="persistent_volumes_claim_ref"></a>
The `claim_ref` block supports:

* `api_version` - The API version.

* `field_path` - The field path.

* `kind` - The kind.

* `name` - The name.

* `namespace` - The namespace.

* `resource_version` - The resource version.

* `uid` - The uid.

<a name="persistent_volumes_csi"></a>
The `csi` block supports:

* `controller_expand_secret_ref` - The controller expand secret reference.
  The [controller_expand_secret_ref](#secret_ref) structure is documented below.

* `controller_publish_secret_ref` - The controller publish secret reference.
  The [controller_publish_secret_ref](#secret_ref) structure is documented below.

* `driver` - The driver.

* `fs_type` - The fs type.

* `node_expand_secret_ref` - The node expand secret reference.
  The [node_expand_secret_ref](#secret_ref) structure is documented below.

* `node_publish_secret_ref` - The node publish secret reference.
  The [node_publish_secret_ref](#secret_ref) structure is documented below.

* `node_stage_secret_ref` - The node stage secret reference.
  The [node_stage_secret_ref](#secret_ref) structure is documented below.

* `read_only` - Whether to read only.

* `volume_attributes` - The volume attributes.

* `volume_handle` - The volume handle.

<a name="secret_ref"></a>
The `controller_expand_secret_ref`, `controller_publish_secret_ref`, `node_expand_secret_ref`,
`node_publish_secret_ref`, `node_stage_secret_ref` block supports:

* `name` - The name.

* `namespace` - The namespace.

<a name="persistent_volumes_node_affinity"></a>
The `node_affinity` block supports:

* `required` - The required.
  The [required](#persistent_volumes_node_affinity_required) structure is documented below.

<a name="persistent_volumes_node_affinity_required"></a>
The `required` block supports:

* `node_selector_terms` - The node selector terms.
  The [node_selector_terms](#required_node_selector_terms) structure is documented below.

<a name="required_node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - The match expressions.
  The [match_expressions](#required_node_selector_terms_match_expressions) structure is documented below.

<a name="required_node_selector_terms_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key.

* `operator` - The operator.

* `values` - The values.

<a name="persistent_volumes_status"></a>
The `status` block supports:

* `message` - The message.

* `phase` - The phase.

* `reason` - The reason.
