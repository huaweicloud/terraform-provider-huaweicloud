---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_persistent_volume_claim"
description: |-
  Manages a CCI persistent volume claim resource within HuaweiCloud.
---

# huaweicloud_cciv2_persistent_volume_claim

Manages a CCI persistent volume claim resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_persistent_volume_claim" "my-pvc" {
  name      = "my-pvc-obs"
  namespace = "default"

  annotations = {
    "everest.io/obs-volume-type"       = "STANDARD"
    "csi.storage.k8s.io/fstype"        = "s3fs"
    "everest.io/enterprise-project-id" = "0"
  }

  access_modes = ["ReadWriteMany"]
  resources {
    requests = {
      storage = "1Gi"
    }
  }
  storage_class_name = "csi-obs"

  volume_mode = "Filesystem"
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) The name of the persistent volume claim in the namespace.

* `namespace` - (Required, String, NonUpdatable) The name of the namespace.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the annotations of the persistent volume claim.

* `labels` - (Optional, Map, NonUpdatable) Specifies the labels of the persistent volume claim.

* `access_modes` - (Optional, List, NonUpdatable) Specifies the access_modes of the persistent volume claim.
  
* `resources` - (Optional, List, NonUpdatable) Specifies the access_modes of the persistent volume claim.
  The [resources](#resources) structure is documented below.

* `selector` - (Optional, List) Specifies the selector of the persistent volume claim.
  The [selector](#selector) structure is documented below.

* `storage_class_name` - (Optional, String, NonUpdatable) Specifies the storage class name of the persistent volume claim.

* `volume_mode` - (Optional, String, NonUpdatable) Specifies the volume mode of the persistent volume claim.

* `valume_name` - (Optional, String, NonUpdatable) Specifies the valume name of the persistent volume claim.

<a name="resources"></a>
The `resources` block supports:

* `limits` - (Optional, Map, NonUpdatable) Specifies the limits expressions of the resources.

* `requests` - (Optional, Map, NonUpdatable) Specifies the requests labels of the resources.

<a name="selector"></a>
The `selector` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions of the selector.
  The [match_expressions](#match_expressions) structure is documented below.

* `match_labels` - (Optional, Map) Specifies the match labels of the selector.

<a name="match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Optional, String) Specifies the key of the match expressions.

* `operator` - (Optional, String) Specifies the operator of the match expressions.

* `values` - (Optional, List) Specifies the values of the match expressions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `annotations` - The annotations of the persistent volume claim.

* `api_version` - The API version of the persistent volume claim.

* `creation_timestamp` - The creation timestamp of the persistent volume claim.

* `finalizers` - The finalizers of the persistent volume claim.

* `kind` - The kind of the persistent volume claim.

* `labels` - The labels of the persistent volume claim.

* `resource_version` - The resource version of the persistent volume claim.

* `status` - The status of the persistent volume claim.

* `uid` - The uid of the persistent volume claim.

## Import

The persistent volume claim can be imported using `namespace` and `name`, e.g.

```bash
$ terraform import huaweicloud_cciv2_persistent_volume_claim.test <namespace/name>
```
