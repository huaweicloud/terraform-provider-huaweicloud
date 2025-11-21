---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_policy_ip_reputation"
description: |-
  Use this data source to get the list of WAF policy IP reputation.
---

# huaweicloud_waf_policy_ip_reputation

Use this data source to get the list of WAF policy IP reputation.

## Example Usage

```hcl
variable "lang" {}
variable "type" {}

data "huaweicloud_waf_policy_ip_reputation" "test" {
  lang = var.lang
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `lang` - (Required, String) Specifies the language. Valid values are **cn** and **en**.

* `type` - (Required, String) Specifies the type. Currently, only **idc** is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ip_reputation_map` - The WAF policy IP reputation list.

  The [ip_reputation_map](ip_reputation_map_struct) structure is documented below.

* `locale` - The threat intelligence control and protection description map.
  + **Dr.Peng**: Dr. Peng Telecom & Media Group is a company that provides services such as internet data centers and
    cloud computing.
  + **Google**: Google, a globally renowned technology company, provides services such as search engines and cloud computing.
  + **Tencent**: Tencent, a well-known Chinese internet company, provides services such as social networking, gaming,
    and finance.
  + **MeiTuan**: Meituan, China's leading e-commerce platform for local services.
  + **Microsoft**: Microsoft Corporation, a globally renowned technology company, provides operating systems, office
    software, and other services.
  + **AliCloud**: Alibaba Cloud, the cloud computing brand under Alibaba Group.
  + **Amazon**: Amazon, a globally renowned e-commerce and cloud computing company.
  + **VNET**: 21Vianet, China's leading telecom-neutral internet infrastructure service provider.
  + **HW**: Huawei, a globally renowned telecommunications technology company.

<a name="ip_reputation_map_struct"></a>
The `ip_reputation_map` block supports:

* `idc` - The content types of threat intelligence control.
