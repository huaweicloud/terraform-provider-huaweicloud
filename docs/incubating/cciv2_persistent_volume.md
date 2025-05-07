---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_persistent_volume"
description: |-
  Manages a CCI persistent volume resource within HuaweiCloud.
---
# huaweicloud_cciv2_persistent_volume

Manages a CCI persistent volume resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "driver" {}
variable "volume_handle" {}

resource "huaweicloud_cciv2_persistent_volume" test {
  name         = var.name
  access_modes = ["ReadWriteMany"]

  capacity = {
    storage = "2Gi"
  }

  csi {
    driver        = var.driver
    volume_handle = var.volume_handle
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI persistent volume.

* `access_modes` - (Required, List) Specifies the access modes of the CCI persistent volume.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI persistent volume.

* `capacity` - (Required, Map) Specifies the capacity of the CCI persistent volume.

* `claim_ref` - (Optional, List) Specifies the claim of the CCI persistent volume.
  The [claim_ref](#claim_ref) structure is documented below.

* `csi` - (Required, List) Specifies the CSI of the CCI persistent volume.
  The [csi](#csi) structure is documented below.

* `labels` - (Optional, Map) Specifies the labels of the CCI persistent volume.

* `mount_options` - (Optional, List) Specifies the mount options of the CCI persistent volume.

* `node_affinity` - (Optional, List) Specifies the node affinity of the CCI persistent volume.
  The [node_affinity](#node_affinity) structure is documented below.

* `reclaim_policy` - (Optional, String) Specifies the reclaim policy of the CCI persistent volume.

* `storage_class_name` - (Optional, String) Specifies the storage class name of the CCI persistent volume.

* `volume_mode` - (Optional, String) Specifies the volume mode of the CCI persistent volume.

<a name="claim_ref"></a>
The `claim_ref` block supports:

* `api_version` - (Optional, String) Specifies the api version of the claim.

* `field_path` - (Optional, String) Specifies the field path of the claim.

* `kind` - (Optional, String) Specifies the kind of the claim.

* `name` - (Optional, String) Specifies the name of the claim.

* `namespace` - (Optional, String) Specifies the namespace of the claim.

* `resource_version` - (Optional, String) Specifies the resource version of the claim.

* `uid` - (Optional, String) Specifies the uid of the claim.

<a name="csi"></a>
The `csi` block supports:

* `driver` - (Required, String, NonUpdatable) Specifies the driver of the CSI.

* `volume_handle` - (Required, String, NonUpdatable) Specifies the volume handle of the CSI.

* `controller_expand_secret_ref` - (Optional, List) Specifies the controller expand secret of the CSI.
  The [controller_expand_secret_ref](#secret_ref) structure is documented below.

* `controller_publish_secret_ref` - (Optional, List) Specifies the controller publish secret of the CSI.
  The [controller_publish_secret_ref](#secret_ref) structure is documented below.

* `fs_type` - (Optional, String) Specifies the FS type of the CSI.

* `node_expand_secret_ref` - (Optional, List) Specifies the node expand secret of the CSI.
  The [node_expand_secret_ref](#secret_ref) structure is documented below.

* `node_publish_secret_ref` - (Optional, List) Specifies the node publish secret of the CSI.
  The [node_publish_secret_ref](#secret_ref) structure is documented below.

* `node_stage_secret_ref` - (Optional, List) Specifies the node stage secret of the CSI.
  The [node_stage_secret_ref](#secret_ref) structure is documented below.

* `read_only` - (Optional, Bool) Specifies whether to read only.

* `volume_attributes` - (Optional, Map) Specifies the volume attributes of the CSI.

<a name="secret_ref"></a>
The `controller_expand_secret_ref`, `controller_publish_secret_ref`, `node_expand_secret_ref`,
`node_publish_secret_ref`, `node_stage_secret_ref` block supports:

* `name` - (Optional, String) Specifies the name of the secret resource.

* `namespace` - (Optional, String) Specifies the namespace of the secret resource.

<a name="node_affinity"></a>
The `node_affinity` block supports:

* `required` - (Optional, List) Specifies the required field of the volume node affinity.
  The [required](#node_affinity_required) structure is documented below.

<a name="node_affinity_required"></a>
The `required` block supports:

* `node_selector_terms` - (Required, List) Specifies the node selector terms.
  The [node_selector_terms](#node_selector_terms) structure is documented below.

<a name="node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions.
  The [match_expressions](#match_expressions) structure is documented below.

<a name="match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Required, String) Specifies the key of the match expressions.

* `operator` - (Required, String) Specifies the operator of the match expressions.

* `values` - (Optional, List) Specifies the values of the match expressions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI persistent volume.

* `creation_timestamp` - The creation timestamp of the CCI persistent volume.

* `finalizers` - The finalizers of the CCI persistent volume.

* `kind` - The kind of the CCI persistent volume.

* `resource_version` - The resource version of the CCI persistent volume.

* `status` - The status of the CCI persistent volume.
  The [status](#status) structure is documented below.

* `uid` - The uid of the CCI persistent volume.

<a name="status"></a>
The `status` block supports:

* `message` - The message.

* `phase` - The phase.

* `reason` - The reason.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The CCI persistent volume can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_cciv2_persistent_volume.test <name>
```
