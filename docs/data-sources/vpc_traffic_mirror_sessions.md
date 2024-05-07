---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_vpc_traffic_mirror_sessions

Use this data source to get the traffic mirror sessions.

## Example Usage

```hcl
data "huaweicloud_vpc_traffic_mirror_sessions" "test1" {
  name = "mirror-session-a6b5"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `traffic_mirror_session_id` - (Optional, String) Specifies the traffic mirror session ID used to query.

* `name` - (Optional, String) Specifies the traffic mirror session name used to query.

* `traffic_mirror_filter_id` - (Optional, String) Specifies the traffic mirror filter ID used in the session.

* `traffic_mirror_target_id` - (Optional, String) Specifies the traffic mirror target ID.

* `traffic_mirror_target_type` - (Optional, String) Specifies the mirror target type. The value can be:
  + **eni**: elastic network interface;
  + **elb**: private network load balancer;

* `priority` - (Optional, String) Specifies the mirror session priority. The value range is **1-32766**.
  A smaller value indicates a higher priority.

* `enabled` - (Optional, String) Specifies whether the mirror session is enabled. Defaults to **true**.

* `type` - (Optional, String) Specifies the mirror source type. The value can be **eni**(elastic network interface).

* `virtual_network_id` - (Optional, String) Specifies the VNI, which is used to distinguish mirrored traffic of different
  sessions. The value range is **0-16777215**, defaults to **1**.

* `packet_length` - (Optional, String) Specifies the maximum transmission unit (MTU).
  The value range is **1-1460**, defaults to **96**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `traffic_mirror_sessions` - List of traffic mirror sessions.

  The [traffic_mirror_sessions](#traffic_mirror_sessions_struct) structure is documented below.

<a name="traffic_mirror_sessions_struct"></a>
The `traffic_mirror_sessions` block supports:

* `project_id` - Project ID.

* `id` - Traffic mirror session ID.

* `name` - Traffic mirror session name.

* `description` - Description of a traffic mirror session.

* `type` - Supported mirror source type.

* `enabled` - Whether to enable a mirror session.

* `priority` - Mirror session priority.

* `packet_length` - Maximum transmission unit (MTU).

* `virtual_network_id` - VNI, which is used to distinguish mirrored traffic of different sessions.

* `traffic_mirror_filter_id` - Traffic mirror filter ID.

* `traffic_mirror_sources` - Mirror source IDs. An elastic network interface can be used as a mirror source.
  Each mirror session can have up to 10 mirror sources by default.

* `traffic_mirror_target_id` - Mirror target ID.

* `traffic_mirror_target_type` - Mirror target type.

* `created_at` - Time when a traffic mirror session is created.

* `updated_at` - Time when the traffic mirror session is updated.
