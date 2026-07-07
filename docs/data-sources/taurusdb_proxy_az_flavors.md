---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_proxy_az_flavors"
description: |
  Use this data source to query available proxy flavors of a TaurusDB instance by AZ codes within HuaweiCloud.
---

# huaweicloud_taurusdb_proxy_az_flavors

Use this data source to query available proxy flavors of a TaurusDB instance by AZ codes within HuaweiCloud.

## Example Usage

```hcl
variable "az_codes" {}

data "huaweicloud_taurusdb_proxy_az_flavors" "test" {
  az_codes           = var.az_codes
  proxy_engine_name  = "taurusproxy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the resource. If omitted, the provider-level region
  will be used.

* `az_codes` - (Required, String) Specifies the AZ codes for querying proxy flavors.
  Multiple AZ codes should be concatenated with commas (`,`), e.g., **cn-north-4a,cn-north-4b**.

* `proxy_engine_name` - (Required, String) Specifies the proxy engine name. The valid value is **taurusproxy**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `proxy_flavor_groups` - Indicates the list of flavor groups.
  The [proxy_flavor_groups](#proxy_flavor_groups_struct) structure is documented below.

<a name="proxy_flavor_groups_struct"></a>
The `proxy_flavor_groups` block contains:

* `group_type` - Indicates the group type. The value can be **arm** or **x86**.

* `proxy_flavors` - Indicates the list of flavors.
  The [proxy_flavors](#proxy_flavors_struct) structure is documented below.

<a name="proxy_flavors_struct"></a>
The `proxy_flavors` block contains:

* `id` - Indicates the ID of the proxy flavor.

* `spec_code` - Indicates the specification code.

* `vcpus` - Indicates the number of vCPUs.

* `ram` - Indicates the memory size.

* `db_type` - Indicates the database type.

* `az_status` - Indicates the AZ status.

* `supported_ipv6` - Indicates whether IPv6 is supported.
