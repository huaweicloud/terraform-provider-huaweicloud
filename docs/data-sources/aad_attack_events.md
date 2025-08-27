---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_attack_events"
description: |-
  Use this data source to get the list of Advanced Anti-DDos attack events within HuaweiCloud.
---

# huaweicloud_aad_attack_events

Use this data source to get the list of Advanced Anti-DDos attack events within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_attack_events" "test" {
  recent        = "today"
  overseas_type = "0"
}
```

## Argument Reference

The following arguments are supported:

* `domains` - (Optional, String) Specifies the domain name. If not specified, all domains will be included.

* `start_time` - (Optional, String) Specifies the start time.

* `end_time` - (Optional, String) Specifies the end time.

* `recent` - (Optional, String) Specifies the recent time period.  
  The valid values are as follows:
  + **yesterday**
  + **today**
  + **3days**
  + **1week**
  + **1month**

  Choose between this parameter and the timestamp parameter(`start_time` and `end_time`).

* `overseas_type` - (Optional, String) Specifies the instance type.  
  The valid values are as follows:
  + **0**: Mainland China.
  + **1**: Overseas.

* `sip` - (Optional, String) Specifies the attack source IP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `events` - The attack events list.  
  The [events](#events_struct) structure is documented below.

<a name="events_struct"></a>
The `events` block supports:

* `id` - The event ID.

* `domain` - The attack target domain.

* `time` - The attack time.

* `sip` - The attack source IP.

* `action` - The defense action.

* `url` - The attack URL.

* `type` - The attack type.

* `backend` - The current backend information.  
  The [backend](#backend_struct) structure is documented below.

<a name="backend_struct"></a>
The `backend` block supports:

* `protocol` - The current backend protocol.

* `port` - The current backend port.

* `host` - The current backend host value.
