---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_quotas"
description: |-
  Use this data source to get the list of HSS quotas within HuaweiCloud.
---

# huaweicloud_hss_quotas

Use this data source to get the list of HSS quotas within HuaweiCloud.

## Example Usage

```hcl
variable quota_id {}

data "huaweicloud_hss_quotas" "test" {
  quota_id = var.quota_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS quotas.
  If omitted, the provider-level region will be used.

* `category` - (Optional, String) Specifies the category of the quotas to be queried.
  The valid values are as follows:
  + **host_resource**: Host protection quota.
  + **container_resource**: Container protection quota.

  If omitted, return all quotas for host resource.  
  If set to **container_resource**, return all quotas with version **hss.version.container.enterprise**.

* `version` - (Optional, String) Specifies the version of the quotas to be queried.
  The valid values are as follows:
  + **hss.version.basic**: Basic version.
  + **hss.version.advanced**: Professional version.
  + **hss.version.enterprise**: Enterprise version.
  + **hss.version.premium**: Ultimate version.
  + **hss.version.wtp**: Web page tamper prevention version.

* `status` - (Optional, String) Specifies the status of the quotas to be queried.
  The value can be **normal**, **expired**, or **freeze**.

* `used_status` - (Optional, String) Specifies the usage status of the quotas to be queried.
  The value can be **idle** or **used**.

* `host_name` - (Optional, String) Specifies the host name for the quota binding to be queried.

* `quota_id` - (Optional, String) Specifies the ID of the quota to be queried.

* `charging_mode` - (Optional, String) Specifies the charging mode of the quotas to be queried.
  The valid values are as follows:
  + **prePaid**: The yearly/monthly billing mode.
  + **postPaid**: The pay-per-use billing mode.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the quotas belong.
  For enterprise users, if omitted, will query the quotas under all enterprise projects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `quotas` - All quotas that match the filter parameters.  
  The [quotas](#hss_quotas) structure is documented below.

<a name="hss_quotas"></a>
The `quotas` block supports:

* `id` - The ID of quota.

* `version` - The version of quota.

* `status` - The status of quota.

* `used_status` - The usage status of quota.

* `host_id` - The host ID for quota binding.

* `host_name` - The host name for quota binding.

* `charging_mode` - The charging mode of quota.

* `expire_time` - The expiration time of quota, in RFC3339 format. This field is valid when the quota is a trial quota.

* `shared_quota` - Is it a shared quota. The value can be **shared** or **unshared**.

* `enterprise_project_id` - The enterprise project ID to which the quota belongs.

* `enterprise_project_name` - The enterprise project name to which the quota belongs.

* `tags` - The key/value pairs to associate with the HSS quota.
