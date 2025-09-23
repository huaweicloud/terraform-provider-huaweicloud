---
subcategory: Cloud Connect (CC)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth"
description: ""
---

# huaweicloud_cc_global_connection_bandwidth

Manages a CC global connection bandwidth within HuaweiCloud.

## Example Usage

```hcl
variable "gcb_name" {}

resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name        = var.gcb_name
  type        = "Region"
  bordercross = false
  charge_mode = "bwd"
  size        = 5
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the GCB name.

* `type` - (Required, String, ForceNew) Specifies the GCB type.
  
  Valid values are as follows:
  + **TrsArea**: Cross geographic region.
  + **Area**: Geographic region.
  + **SubArea**: Homezones region.
  + **Region**: Multi-city region.

  Changing this creates a new resource.

* `bordercross` - (Required, Bool, ForceNew) Specifies whether the GCB involves traveling from Chinese mainland to other
  countries. Changing this creates a new resource.

* `charge_mode` - (Required, String) Specifies the GCB charge mode.

  Valid values are as follows:
  + **bwd**: Billed by bandwidth.
  + **95**: Billed by 95th percentile bandwidth.

  Only support changing **bwd** to **95**.

* `size` - (Required, Int) Specifies the GCB size. If `charge_mode` is **bwd**, value ranges from `2` to `300` Mbit/s. If
  `charge_mode` is **95**, value ranges from `100` to `300` Mbit/s.

* `sla_level` - (Optional, String) Specifies the network level. From high to low, divided into **Pt**(platinum),
  **Au**(gold), and **Ag**(silver). The default is **Au**.

* `local_area` - (Optional, String, ForceNew) Specifies the local access point. The valid length is limited between `1`
  to `64`, Only Chinese and English letters, digits, hyphens (-), underscores (_) and dots (.) are allowed. If `type` is
  **Region**, it is **optional**, otherwise it is **Required**.
  Changing this creates a new resource.

* `remote_area` - (Optional, String, ForceNew) Specifies the remote access point. The valid length is limited between `1`
  to `64`, Only Chinese and English letters, digits, hyphens (-), underscores (_) and dots (.) are allowed. If `type` is
  **Region**, it is **optional**, otherwise it is **Required**.
  Changing this creates a new resource.

* `spec_code_id` - (Optional, String) Specifies the line specification code UUID.

* `binding_service` - (Optional, String) Specifies whether to limit the GCB only bind with specific instance. Default is
  **ALL**.

  Valid values are as follows:
  + **CC**: Cloud Connection.
  + **GEIP**: Global EIP.
  + **GCN**: Central Network.
  + **GSN**: Site Network.
  + **ALL**: All.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the GCB belongs.
  Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of GCB. Not support angle brackets (<>).

* `tags` - (Optional, Map) Specifies the tags of GCB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `enable_share` - Indicates the GCB whether to support binding multiple instances.

* `frozen` - Indicates the GCB is frozen or not.

* `instances` - The instances which the GCB binding with.
  The [instances](#attrblock--instances) structure is documented below.

* `created_at` - The create time of GCB.

* `updated_at` - The update time of GCB.

<a name="attrblock--instances"></a>
The `instances` block supports:

* `id` - The instance ID.

* `region` - The region of the instance.

* `type` - The type of the instance.

## Import

The global connection bandwidth can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cc_global_connection_bandwidth.test <id>
```
