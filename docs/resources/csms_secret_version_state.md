---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_secret_version_state"
description: |
  Manages a CSMS secret version state resource within HuaweiCloud.
---

# huaweicloud_csms_secret_version_state

Manages a CSMS secret version state resource within HuaweiCloud.

-> A secret supports a maximum of `12` secret version states, each secret version state can identify only one
  secret version.
  <br>If you add a secret version state in use to a new secret version, the secret version state will be
  automatically removed from the old secret version.
  <br>**SYSCURRENT** and **SYSPREVIOUS** are built-in states, not support deletion.

## Example Usage

```hcl
variable "secret_name" {}
variable "name" {}
variable "version_id" {}

resource "huaweicloud_csms_secret_version_state" "test" {
  secret_name = var.secret_name
  name        = var.name
  version_id  = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CSMS secret version state.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `secret_name` - (Required, String, ForceNew) Specifies the name of the secret to which the secret version state
  belongs. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the secret version state.
  Changing this parameter will create a new secret version.
  Only letters, digits, underscores(_) and hyphens(-) are allowed.
  The valid length is limited from `1` to `64` characters.

* `version_id` - (Required, String) Specifies the ID of the secret version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `name`.

* `updated_at` - The last update time of the secret version state, in RFC3339 format.

## Import

The secret version state can be imported using the related `secret_name` and their `id`, separated by a slash (/), e.g.

```bash
terraform import huaweicloud_csms_secret_version_state.test <secret_name>/<id>
```
