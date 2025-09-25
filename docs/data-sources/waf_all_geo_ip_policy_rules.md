---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_geo_ip_policy_rules"
description: |-
  Use this data source to get list of the WAF geo IP policy rules under all policies.
---

# huaweicloud_waf_all_geo_ip_policy_rules

Use this data source to get list of the WAF geo IP policy rules under all policies.

## Example Usage

```hcl
data "huaweicloud_waf_all_geo_ip_policy_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policyids` - (Optional, String) Specifies the ID of the policy.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of the WAF geo IP policy rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the rule.

* `policyid` - The ID of the policy.

* `name` - The name of the rule.

* `geoip` - The geo location code blocked by the rule.
  The valid values are as follows:
  + **CA**: Canada
  + **US**: United States
  + **AU**: Australia
  + **IN**: India
  + **JP**: Japan
  + **UK**: United Kingdom
  + **FR**: France
  + **DE**: Germany
  + **BR**: Brazil
  + **Thailand**: Thailand
  + **Singapore**: Singapore
  + **South Africa**: South Africa
  + **Mexico**: Mexico
  + **Peru**: Peru
  + **Indonesia**: Indonesia
  + **GD**: Guangdong (China)
  + **FJ**: Fujian (China)
  + **JL**: Jilin (China)
  + **LN**: Liaoning (China)
  + **TW**: Taiwan (China)
  + **GZ**: Guizhou (China)
  + **AH**: Anhui (China)
  + **HL**: Heilongjiang (China)
  + **HA**: Henan (China)
  + **SC**: Sichuan (China)
  + **HE**: Hebei (China)
  + **YN**: Yunnan (China)
  + **HB**: Hubei (China)
  + **HI**: Hainan (China)
  + **QH**: Qinghai (China)
  + **HN**: Hunan (China)
  + **JX**: Jiangxi (China)
  + **SX**: Shanxi (China)
  + **SN**: Shaanxi (China)
  + **ZJ**: Zhejiang (China)
  + **GS**: Gansu (China)
  + **JS**: Jiangsu (China)
  + **SD**: Shandong (China)
  + **BJ**: Beijing (China)
  + **SH**: Shanghai (China)
  + **TJ**: Tianjin (China)
  + **CQ**: Chongqing (China)
  + **MO**: Macao (China)
  + **HK**: Hong Kong (China)
  + **NX**: Ningxia (China)
  + **GX**: Guangxi (China)
  + **XJ**: Xinjiang (China)
  + **XZ**: Tibet (China)
  + **NM**: Inner Mongolia (China)

* `white` - The protection action of the rule.
  The valid values are as follows:
  + `0`: Block
  + `1`: Allow
  + `2`: Log only

* `status` - The status of the rule.
  The valid values are as follows:
  + `0`: Disabled
  + `1`: Enabled

* `timestamp` - The creation time of the rule, in milliseconds.
