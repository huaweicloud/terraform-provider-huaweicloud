---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_resource_quotas"
description: |-
  Use this data source to get the resource quotas of HSS within HuaweiCloud.
---

# huaweicloud_hss_resource_quotas

Use this data source to get the resource quotas of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_resource_quotas" "test" {
  version = "hss.version.basic"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `version` - (Optional, String) Specifies the HSS version. The valid values are as follows:
  + **hss.version.null**: No version.
  + **hss.version.basic**: Basic edition.
  + **hss.version.advanced**: Professional edition.
  + **hss.version.enterprise**: Enterprise edition.
  + **hss.version.premium**: Premium edition.
  + **hss.version.wtp**: Web tamper protection edition.
  + **hss.version.container.enterprise**: Container edition.

* `charging_mode` - (Optional, String) Specifies the billing mode. The valid values are as follows:
  + **packet_cycle**: Yearly/Monthly subscription.
  + **on_demand**: Pay-per-use.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The quota information list.
  The [data_list](#quota_info_structure) structure is documented below.

<a name="quota_info_structure"></a>
The `data_list` block supports:

* `version` - The HSS version The valid values are as follows:
  + **hss.version.null**: No version.
  + **hss.version.basic**: Basic edition.
  + **hss.version.advanced**: Professional edition.
  + **hss.version.enterprise**: Enterprise edition.
  + **hss.version.premium**: Premium edition.
  + **hss.version.wtp**: Web tamper protection edition.
  + **hss.version.container.enterprise**: Container edition.

* `total_num` - The total quota number.

* `used_num` - The used quota number.

* `available_num` - The available quota number.

* `available_resources_list` - The list of available resources.
  The [available_resources_list](#available_resource_structure) structure is documented below.

<a name="available_resource_structure"></a>
The `available_resources_list` block supports:

* `resource_id` - The resource ID.

* `current_time` - The current time.

* `shared_quota` - Whether the quota is shared. The valid values are:
  + **shared**: The quota is shared.
  + **unshared**: The quota is not shared.
  