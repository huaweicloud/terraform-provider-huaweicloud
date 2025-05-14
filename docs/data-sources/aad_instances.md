---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_instances"
description: |-
  Use this data source to get the list of Advanced Anti-DDos instances within HuaweiCloud.
---

# huaweicloud_aad_instances

Use this data source to get the list of Advanced Anti-DDos instances within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_instances" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `items` - The list of the AAD instances.
  The [items](#items) structure is documented below.

<a name="items"></a>
The `items` block supports:

* `instance_id` - The AAD instance ID.

* `instance_name` - The name of the AAD instance.

* `ips` - The list of the AAD instance IPs.
  The [ips](#ips) structure is documented below.

* `expire_time` - The expiration time of the AAD instance.

* `service_bandwidth` - The service bandwidth of the AAD instance.

* `instance_status` - The AAD instance status.

* `enterprise_project_id` - The enterprise project ID of the AAD instance.

* `overseas_type` - The AAD instance type, `0`-mainland China, `1`-overseas.

<a name="ips"></a>
The `ips` block supports:

* `ip_id` - The IP ID of the AAD instance.

* `ip` - The IP of the AAD instance.

* `basic_bandwidth` - The basic bandwidth of the AAD instance.

* `elastic_bandwidth` - The elastic bandwidth of the AAD instance.

* `ip_status` - The IP status of the AAD instance.
