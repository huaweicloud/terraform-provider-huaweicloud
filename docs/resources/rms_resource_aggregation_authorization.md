---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregation_authorization"
description: ""
---

# huaweicloud_rms_resource_aggregation_authorization

Manages a RMS aggregation authorization resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_account" {}

resource "huaweicloud_rms_resource_aggregation_authorization" "test" {
  account_id = var.source_account
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, String, ForceNew) Specifies the ID of the resource aggregation account to be authorized.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map)  Specifies the key/value pairs to associate with the aggregation authorization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the authorization identifier of the resource aggregation account.

* `created_at` - Indicates the time when the resource aggregation account was authorized.

## Import

The aggregation authorization can be imported using the `account_id`, e.g.

```bash
$ terraform import huaweicloud_rms_resource_aggregation_authorization.test 036a12ef8327c4194346684fdbe0b37e
```
