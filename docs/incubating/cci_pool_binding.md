---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cci_pool_binding"
description: |-
  Manages a CCI pool binding resource within HuaweiCloud.
---

# huaweicloud_cci_pool_binding

Manages a CCI pool binding resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}
variable "id" {}

resource "huaweicloud_cci_pool_binding" "test" {
  namespace = var.namespace
  name      = var.name

  pool_ref {
    id = var.id
  }

  target_ref {
    group = "cci/v2"
    kind  = "Service"
    name  = "test-service"
    port  = 1234
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies a unique name within a namespace.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace.

* `annotations` - (Optional, Map) Specifies the namespace.

* `api_version` - (Optional, String) Specifies the API version. Defaults to **loadbalancer.networking.openvessel.io/v1**.

* `finalizers` - (Optional, List) Specifies the finalizers.

* `generate_name` - (Optional, String) Specifies the generate name.

* `kind` - (Optional, String) Specifies the kind. Defaults to **PoolBinding**.

* `labels` - (Optional, Map) Specifies the labels.

* `owner_references` - (Optional, List) Specifies the owner references.
  The [owner_references](#owner_references) structure is documented below.

* `pool_ref` - (Optional, List) Specifies the pool ref.
  The [pool_ref](#pool_ref) structure is documented below.

* `target_ref` - (Optional, List) Specifies the target ref.
  The [target_ref](#target_ref) structure is documented below.

<a name="owner_references"></a>
The `owner_references` block supports:

* `api_version` - (Required, String) Specifies the API version of the referent.

* `kind` - (Required, String) Specifies the kind of the referent.

* `name` - (Required, String) Specifies the name of the referent.

* `uid` - (Required, String) Specifies the uid of the referent.

* `block_owner_deletion` - (Optional, Bool) Specifies whether it can be deleted the block owner reference.

* `controller` - (Optional, Bool) Specifies whether it is controllered.

<a name="pool_ref"></a>
The `pool_ref` block supports:

* `id` - (Optional, String) Specifies the ID of the ELB backend server group.

<a name="target_ref"></a>
The `target_ref` block supports:

* `name` - (Required, String) Specifies the name of the target reference.

* `group` - (Optional, String) Specifies the group of the target reference.

* `kind` - (Optional, String) Specifies the kind of the target reference.

* `namespace` - (Optional, String) Specifies the namespace of the target reference.

* `port` - (Optional, Int) Specifies the port of the target reference.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creation_timestamp` - The creation time.

* `deletion_grace_period_seconds` - The deletion grace period seconds.

* `deletion_timestamp` - The deletion time.

* `generation` - The generation.

* `resource_version` - The resource version.

* `uid` - The uid.

## Import

The CCI pool binding can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cci_pool_binding.test <namespace>/<name>
```
