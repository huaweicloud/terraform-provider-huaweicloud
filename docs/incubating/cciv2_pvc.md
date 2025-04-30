---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_pvc"
description: |-
  Manages a CCI persistent volume claim resource within HuaweiCloud.
---

# huaweicloud_cciv2_pvc

Manages a CCI persistent volume claim resource within HuaweiCloud.

## Example Usage

<!-- please add the usage of huaweicloud_cciv2_pvc -->
```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) The name of the persistent volume claim in the namespace.

* `namespace` - (Required, String, NonUpdatable) The name of the namespace.

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

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

The xxx can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_pvc.test <id>
```
