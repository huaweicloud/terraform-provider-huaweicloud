---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidths"
description: ""
---

# huaweicloud_cc_global_connection_bandwidths

Use this data source to get the list of CC global connection bandwidths.

## Example Usage

```hcl
variable "gcb_id" {}
variable "name" {}

data "huaweicloud_cc_global_connection_bandwidths" "test" {
  gcb_id = var.gcb_id
  name   = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.

* `name` - (Optional, String) Specifies the resource name.

* `gcb_id` - (Optional, String) Specifies the resource ID.

* `size` - (Optional, Int) Specifies the bandwidth range.
  Bandwidth range: `2` Mbit/s to `300` Mbit/s.

* `admin_state` - (Optional, String) Specifies the status of the global connection bandwidth.
  Value options are as follows:
  + **NORMAL**: The global connection bandwidth is available.
  + **FREEZED**: The global connection bandwidth is frozen.

* `type` - (Optional, String) Specifies the type of a global connection bandwidth.
  Value options are as follows:
  + **TrsArea**: cross-geographic-region bandwidth.
  + **Area**: geographic-region bandwidth.
  + **SubArea**: region bandwidth.
  + **Region** : multi-city bandwidth.

* `instance_id` - (Optional, String) Specifies the bound instance ID.

* `instance_type` - (Optional, String) Specifies the instance type.

* `binding_service` - (Optional, String) Specifies the binding service.
  Value options are as follows:
  + **Cloud Connect**: cloud connection.
  + **GEIP**: Global EIP.
  + **GCN**: central network.
  + **GSN**: site network.

* `charge_mode` - (Optional, String) Specifies the billing option.
  Value options are as follows:
  + **bwd**: billing by bandwidth
  + **95**: standard 95th percentile bandwidth billing

* `tags` - (Optional, Map) Specifies tags.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `globalconnection_bandwidths` - Response body for querying the global connection bandwidth list.
  The [globalconnection_bandwidths](#GlobalConnectionBandwidths) structure is documented below.

<a name="GlobalConnectionBandwidths"></a>
The `globalconnection_bandwidths` block supports:

* `id` - Resource ID.

* `name` - Resource name.

* `description` - Resource description.

* `domain_id` - ID of the account that the resource belongs to.

* `type` - Type of a global connection bandwidth.

* `bordercross` - Whether the global connection bandwidth is used for cross-border communications.

* `binding_service` - Binding service.

* `enterprise_project_id` - ID of the enterprise project that the global connection bandwidth belongs to.

* `charge_mode` - Billing option.

* `size` - Range of a global connection bandwidth, in Mbit/s.

* `sla_level` - Class of a global connection bandwidth.

* `local_site_code` - Code of the local access point.

* `remote_site_code` - Code of the remote access point.

* `admin_state` - Global connection bandwidth status.

* `remote_area` - Name of a remote access point.

* `local_area` - Name of a local access point.

* `frozen` - Whether a global connection bandwidth is frozen.

* `spec_code_id` - UUID of a line specification code.

* `tags` - Resource tags.

* `created_at` - Time when the resource was created.

* `updated_at` - Time when the resource was updated.

* `enable_share` - Whether a global connection bandwidth can be used by multiple instances.

* `instances` - The list of instances that the global connection bandwidth is bound to.
  The [instances](#GlobalConnectionBandwidth_Instances) structure is documented below.

<a name="GlobalConnectionBandwidth_Instances"></a>
The `instances` block supports:

* `id` - Bound instance ID.

* `project_id` - Project ID of the bound instance.

* `region_id` - Region of the bound instance. The default value is **global** for global services.

* `type` - Bound instance type.
