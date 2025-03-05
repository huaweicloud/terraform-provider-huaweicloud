---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_namespace"
description: ""
---

# huaweicloud_cciv2_namespace

Manages a CCI namespace resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace_name" {}

resource "huaweicloud_cci_namespace" "test" {
  name        = var.namespace_name

  annotations = {
    "namespace.kubernetes.io/flavor"            = "gpu-accelerated"
    "network.cci.io/warm-pool-size"             = "10"
    "network.cci.io/warm-pool-recycle-interval" = "24"
    "network.cci.io/ready-before-pod-run"       = "vpc-network-ready"
  }

  labels = {
    "rbac.authorization.cci.io/enable-k8s-rbac" = "true",
    "sys_enterprise_project_id" = "0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the namespace.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the annotations of the namespace.

* `labels` - (Optional, Map, NonUpdatable) Specifies the labels of the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the namespace.

* `kind` - The kind of the namespace.

* `cluster_name` - The cluster name of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `deletion_grace_period_seconds` - The deletion grace period seconds of the namespace.

* `deletion_timestamp` - The deletion timestamp of the namespace.

* `finalizers` - The finalizers of the namespace.

* `generate_name` - The generate name of the namespace.

* `generation` - The generation of the namespace.

* `managed_fields` - The managed fields of the namespace.
  The [managed_fields](#attrblock--managed_fields) structure is documented below.

* `owner_references` - The owner references of the namespace.
  The [owner_references](#attrblock--owner_references) structure is documented below.

* `resource_version` - The resource version of the namespace.

* `self_link` - The self link of the namespace.

* `status` - The status of the namespace.

* `uid` - The uid of the namespace.

<a name="attrblock--managed_fields"></a>
The `managed_fields` block supports:

* `api_version` - The API version of the managed fields.

* `fields_type` - The fields type of the managed fields.

* `fields_v1` - The fields v1 of the managed fields.

* `manager` - The manager of the managed fields.

* `operation` - The operation of the managed fields.

* `time` - The time of the managed fields.

<a name="attrblock--owner_references"></a>
The `owner_references` block supports:

* `api_version` - The API version of the owner references.

* `block_owner_deletion` - The block owner deletion of the owner references.

* `controller` - The controller of the owner references.

* `kind` - The kind of the owner references.

* `name` - The name of the owner references.

* `uid` - The uid of the owner references.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

The xxx can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_namespace.test <id>
```
