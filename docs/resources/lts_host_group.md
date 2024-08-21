---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_host_group"
description: |-
  Manages an LTS host group resource within HuaweiCloud.
---

# huaweicloud_lts_host_group

Manages an LTS host group resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_name" {}
variable "host_ids" {
  type = list(string)
}

resource "huaweicloud_lts_host_group" "test" {
  name     = var.group_name
  type     = "linux"
  host_ids = var.host_ids

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the host group.

* `type` - (Required, String, ForceNew) Specifies the type of the host.
  The value can be **linux** and **windows**.

  Changing this parameter will create a new resource.

* `host_ids` - (Optional, List) Specifies the ID list of hosts to join the host group.

* `agent_access_type` - (Optional, String) Specifies the type of the host group.
  The default value is **IP**.  
  The valid values are as follows:
  + **IP**
  + **LABEL**

* `labels` - (Optional, List) Specifies the custom label list of the host group.
  This parameter is required when `agent_access_type` is set to **LABEL**.

  -> Currently, a maximum of `10` labels can be created.

* `tags` - (Optional, Map) Specifies the key/value to attach to the host group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

## Import

The host group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_host_group.test <id>
```
