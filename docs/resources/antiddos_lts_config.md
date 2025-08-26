---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_lts_config"
description: |-
  Manages Anti-DDoS LTS configuration resource within HuaweiCloud.
---

# huaweicloud_antiddos_lts_config

Manages Anti-DDoS LTS configuration resource within HuaweiCloud.

-> Destroying the resource will disable the LTS configuration.

## Example Usage

```hcl
variable "lts_group_id" {}
variable "lts_attack_stream_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_antiddos_lts_config" "test" {
  lts_group_id          = var.lts_group_id
  lts_attack_stream_id  = var.lts_attack_stream_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `lts_group_id` - (Required, String) Specifies the LTS group ID.

* `lts_attack_stream_id` - (Required, String) Specifies the LTS attack stream ID.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `enterprise_project_id`).

## Import

Anti-DDoS LTS configuration can be imported using `id`. e.g.

```bash
$ terraform import huaweicloud_antiddos_lts_config.test <id>
```
