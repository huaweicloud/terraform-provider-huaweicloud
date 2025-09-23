---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_domains"
description: |-
  Use this data source to get list of the domains, including cloud domains and dedicated domains.
---

# huaweicloud_waf_all_domains

Use this data source to get list of the domains, including cloud domains and dedicated domains.

## Example Usage

```hcl
data "huaweicloud_waf_all_domains" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `hostname` - (Optional, String) Specifies the name of the domain.

* `policyname` - (Optional, String) Specifies the name of the policy.

* `protect_status` - (Optional, Int) Specifies the protection status of the domain.
  The valid values are as follows:
  + `-1`: The WAF protection is bypassed. Requests of the domain are directly sent to the backend server and do not
  pass through WAF.
  + `0`: The WAF protection is suspended. WAF only forwards requests destined for the domain and does not
  detect attacks.
  + `1`: The WAF protection is enabled. WAF detects attacks based on the policy you configure.

* `waf_type` - (Optional, String) Specifies the WAF mode of the domain.
  The valid values are as follows:
  + **cloud**: The cloud WAF is used to protect the domain.
  + **premium**: The dedicated WAF instance is used to protect the domain.

* `is_https` - (Optional, String) Specifies whether HTTPS is used for the domain.
  The value can be **true** of **false**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The details about the protected domain.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the domain.

* `hostname` - The name of the domain.

* `policyid` - The ID of the policy.

* `description` - The description of the domain.

* `protect_status` - The protection status of the domain.

* `access_status` - The domain name access status.
  The valid values are as follows:
  + `0`: The website traffic has not been routed to WAF.
  + `1`: The website traffic has been routed to WAF.

* `access_code` - The cname prefix.

* `proxy` - Whether a proxy is used for the domain.
  The valid values are as follows:
  + **false**: No proxy is used.
  + **true**: A proxy is used.

* `web_tag` - The website name.

* `paid_type` - The package billing mode. The value can be prePaid or postPaid.
  The valid values are as follows:
  + **prePaid**: Indicates yearly/monthly billing.
  + **postPaid**: Indicates pay-per-use billing.

* `waf_type` - The mode of the domain.

* `region` - The region ID.

* `enterprise_project_id` - The enterprise project ID.

* `timestamp` - The creation time of the domain.

* `flag` - The special identifier, which is used on the console.

  The [flag](#items_flag_struct) structure is documented below.

* `server` - The origin server settings of the domain.

  The [server](#items_server_struct) structure is documented below.

* `access_progress` - The access progress.

  The [access_progress](#items_access_progress_struct) structure is documented below.

* `premium_waf_instances` - The list of dedicated WAF instances.

  The [premium_waf_instances](#items_premium_waf_instances_struct) structure is documented below.

<a name="items_flag_struct"></a>
The `flag` block supports:

* `pci_3ds` - Whether the website passes the PCI 3DS certification check.

* `pci_dss` - Whether the website passed the PCI DSS certification check.

* `cname` - The cname record being used.
  The valid values are as follows:
  + **old**: The old cname record is used.
  + **new**: The new cname record is used.

* `ipv6` - Whether IPv6 protection is supported.

* `is_dual_az` - Whether WAF support Multi-AZ mode.

<a name="items_server_struct"></a>
The `server` block supports:

* `front_protocol` - The protocol used by the client to request access to the origin server.

* `back_protocol` - The protocol used by WAF to forward client requests it received to origin servers.

* `type` - The origin server type.
  The value can be **IPv4** or **IPv6**.

* `port` - The port used by WAF to forward client requests to the origin server.

* `address` - The IP address of origin server requested by the client.

* `weight` - The weight of the origin server.

* `vpc_id` - The VPC ID.
  The value is returned only for domain protected with dedicated instances.

<a name="items_access_progress_struct"></a>
The `access_progress` block supports:

* `status` - The status of the access. The value can be **0** or **1**.
  + `0`: The step has not been finished.
  + `1`: The step has finished.

* `step` - The procedure.
  The valid values are as follows.
  + `1`: Whitelisting the WAF IP addresses.
  + `2`: Testing connectivity.
  + `3`: Modifying DNS records.

<a name="items_premium_waf_instances_struct"></a>
The `premium_waf_instances` block supports:

* `id` - The ID of the dedicated WAF instance.

* `name` - The name of the dedicated WAF instance.

* `accessed` - Whether the domain name is added to the dedicated WAF instance.
