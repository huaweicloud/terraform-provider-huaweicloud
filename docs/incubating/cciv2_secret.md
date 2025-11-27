---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_secret"
description: |-
  Manages a CCI v2 Secret resource within HuaweiCloud.
---

# huaweicloud_cciv2_secret

Manages a CCI v2 Secret resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}
variable "data" {}

resource "huaweicloud_cciv2_secret" "test" {
  namespace = var.namespace
  name      = var.name
  data      = var.data
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace of the CCI Secret.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI Secret.

-> When creating the secret for AOM integration, the name must be `cci-aom-app-secret`.

* `string_data` - (Optional, Map) Specifies string data of the CCI Secret.

* `data` - (Optional, Map) Specifies the data of the CCI Secret.

* `type` - (Optional, String) Specifies the type of the CCI Secret.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `annotations` - The annotations of the CCI Secret.

* `api_version` - The API version of the CCI Secret.

* `creation_timestamp` - The creation timestamp of the CCI Secret.

* `immutable` - The immutable of the CCI Secret.

* `kind` - The kind of the CCI Secret.

* `labels` - The labels of the CCI Secret.

* `resource_version` - The resource version of the CCI Secret.

* `uid` - The uid of the CCI Secret.

## Import

The CCI v2 Secret can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_secret.test <namespace>/<name>
```
