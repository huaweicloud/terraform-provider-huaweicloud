---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_support_binding_global_connection_bandwidths"
description: |-
  Use this data source to get the list of CC global connection bandwidths which support binding.
---

# huaweicloud_cc_support_binding_global_connection_bandwidths

Use this data source to get the list of CC global connection bandwidths which support binding.

## Example Usage

```hcl
variable "binding_service" {}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "test" { 
  binding_service = var.binding_service
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `binding_service` - (Required, String) Specifies the binding service.
  The valid values are as follows:
  + **CC**: cloud connection.
  + **GEIP**: global EIP.
  + **GCN**: central network.
  + **GSN**: site network.

* `local_area` - (Optional, String) Specifies the local access point.
  If the bandwidth type is set to **region**, all multi-city bandwidths that meet the filtering criteria are returned.
  This field is not matched for filtering. For other types, this field is used to match **local_area** of the backbone
  bandwidth.

* `remote_area` - (Optional, String) Specifies the remote access point.
  If the bandwidth type is set to **region**, all multi-city bandwidths that meet the filtering criteria are returned.
  This field is not matched for filtering. For other types, this field is used to match **remote_area** of the backbone
  bandwidth.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `gcb_id` - (Optional, String) Specifies the global connection bandwidth ID.

* `name` - (Optional, String) Specifies the global connection bandwidth name.

* `type` - (Optional, String) Specifies the global connection bandwidth type.
  The valid values are as follows:
  + **TrsArea**: cross-geographic-region bandwidth.
  + **Area**: geographic-region bandwidth.
  + **SubArea**: region bandwidth.
  + **Region**: multi-city bandwidth.

* `charge_mode` - (Optional, String) Specifies the billing option.
  The valid values are as follows:
  + **bwd**: billing by bandwidth.
  + **95**: standard 95th percentile bandwidth billing.

* `size` - (Optional, Int) Specifies the global connection bandwidth size.

* `sla_level` - (Optional, String) Specifies the class of a global connection bandwidth.
  The valid values are as follows:
  + **Pt**: platinum.
  + **Au**: gold.
  + **Ag**: silver.

* `admin_state` - (Optional, String) Specifies the global connection bandwidth status.
  The valid values are as follows:
  + **NORMAL**: The bandwidth is normal.
  + **FREEZED**: The bandwidth is frozen.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `globalconnection_bandwidths` - The global connection bandwidth list.

  The [globalconnection_bandwidths](#globalconnection_bandwidths_struct) structure is documented below.

<a name="globalconnection_bandwidths_struct"></a>
The `globalconnection_bandwidths` block supports:

* `id` - The global connection bandwidth ID.

* `name` - The global connection bandwidth name.

* `description` - The global connection bandwidth description.

* `domain_id` - The ID of the account that the global connection bandwidth belongs to.

* `bordercross` - Whether the global connection bandwidth is used for cross-border communications.

* `type` - The type of a global connection bandwidth.

* `binding_service` - The type of the binding service.

* `enterprise_project_id` - The ID of the enterprise project that the global connection bandwidth belongs to.

* `charge_mode` - The billing option. By default, billing by bandwidth is enabled.

* `size` - The range of a global connection bandwidth, in Mbit/s.

* `sla_level` - The class of a global connection bandwidth.

* `local_site_code` - The code of the local access point.

* `remote_site_code` - The code of the remote access point.

* `frozen` - Whether a global connection bandwidth is frozen.

* `spec_code_id` - The UUID of a line specification code.

* `created_at` - Time when the resource was created.

* `updated_at` - Time when the resource was updated.

* `enable_share` - Whether a global connection bandwidth can be used by multiple instances.

* `local_area` - The name of a local access point.

* `remote_area` - The name of a remote access point.

* `admin_state` - The global connection bandwidth status.
