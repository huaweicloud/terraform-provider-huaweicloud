---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_hosts"
description: |-
  Use this data source to get the list of HSS web tamper hosts within HuaweiCloud.
---

# huaweicloud_hss_webtamper_hosts

Use this data source to get the list of HSS web tamper hosts within HuaweiCloud.

## Example Usage

```hcl
variable host_id {}

data "huaweicloud_hss_webtamper_hosts" "test" {
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS web tamper hosts.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the ID of the web tamper host to be queried.

* `name` - (Optional, String) Specifies the name of the web tamper host to be queried.
  This field will undergo a fuzzy matching query, the query result is for all web tamper hosts whose names contain this
  value.

* `public_ip` - (Optional, String) Specifies the elastic public IP address of the web tamper host to be queried.

* `private_ip` - (Optional, String) Specifies the private IP address of the web tamper host to be queried.

* `group_name` - (Optional, String) Specifies the host group name to which the web tamper hosts belong to be queried.

* `os_type` - (Optional, String) Specifies the operating system type of the web tamper host to be queried.
  The value can be **linux** or **windows**.

* `protect_status` - (Optional, String) Specifies the protection status of the web tamper hosts to be queried.
  The value can be **closed** or **opened**.

* `rasp_protect_status` - (Optional, String) Specifies the dynamic protection status of the web tamper hosts to be
  queried. The value can be **closed** or **opened**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the web tamper hosts
  belong. For enterprise users, if omitted, will query the web tamper hosts under all enterprise projects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `hosts` - All web tamper hosts that match the filter parameters.  
  The [hosts](#hss_webtamper_hosts) structure is documented below.

<a name="hss_webtamper_hosts"></a>
The `hosts` block supports:

* `id` - The ID of the web tamper host.

* `name` - The name of the web tamper host.

* `public_ip` - The elastic public IP address of the web tamper host.

* `private_ip` - The private IP address of the web tamper host.

* `group_name` - The host group name to which the web tamper host belongs.

* `os_bit` - The operating system bits of the web tamper host.

* `os_type` - The operating system type of the web tamper host.

* `protect_status` - The protection status of the web tamper host.

* `rasp_protect_status` - The dynamic protection status of the web tamper host.

* `anti_tampering_times` - The number of defended tampering attacks.

* `detect_tampering_times` - The number of detected tampering attacks.
