---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_instance"
description: |-
  Manages an Identity Center instance resource within HuaweiCloud.
---

# huaweicloud_identitycenter_instance

Manages an Identity Center instance resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_identitycenter_instance" "test" {}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Optional, String) Specifies the alias of the user-defined identity store ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `identity_store_id` - The ID of the identity store.

* `instance_urn` - The urn of the instance.

## Import

The Identity Center instance can be imported using the `instance_id`, e.g.

```bash
$ terraform import huaweicloud_identitycenter_instance.test <instance_id>
```
