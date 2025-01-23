---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ips_rule_mode_change"
description: |-
  Manages a CFW IPS rule mode change resource within HuaweiCloud.
---

# huaweicloud_cfw_ips_rule_mode_change

Manages a CFW IPS rule mode change resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
Deleting this resource will not clear the corresponding request record,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "object_id" {}
variable "ips_ids" {}
variable "status" {}

resource "huaweicloud_cfw_ips_rule_mode_change" "test"{
  object_id = var.object_id
  ips_ids   = var.ips_ids
  status    = var.status
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID.

* `status` - (Required, String, NonUpdatable) Specifies the IPS rule status.
  The valid value can be **OBSERVE**, **ENABLE**, **CLOSE**, **DEFAULT** or **ALL_DEFAULT**.

* `ips_ids` - (Optional, List, NonUpdatable) Specifies the IPS rule ID list.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
