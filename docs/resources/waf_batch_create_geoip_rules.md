---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_geoip_rules"
description: |-
  Manages a resource to batch create WAF geoip rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_geoip_rules

Manages a resource to batch create WAF geoip rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating geoip rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_geoip_rules" "test" {
  name                  = "test_rule"
  geoip                 = "CA"
  white                 = 1
  ip_type               = "v4"
  policy_ids            = var.policy_ids
  enterprise_project_id = var.enterprise_project_id
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the rule.

* `geoip` - (Required, String, NonUpdatable) Specifies the geoip of the rule.
  Valid locations are as follows:
  **CN**: China, **CA**: Canada, **US**: The United States, **AU**: Australia, **IN**: India, **JP**: Japan,
  **UK**: United Kingdom, **FR**: France, **DE**: Germany, **BR**: Brazil, **Thailand**: Thailand,
  **Singapore**: Singapore,**South Africa**: South Africa, **Mexico**: Mexico, **Peru**: Peru, **Indonesia**: Indonesia,
  **GD**: Guangdong, **FJ**: Fujian, **JL**: Jilin, **LN**: Liaoning, **TW**: Taiwan SAR, China,**GZ**: Guizhou,
  **AH**: Anhui, **HL**: Heilongjiang, **HA**: Henan, **SC**: Sichuan, **HE**: Hebei, **YN**: Yunnan, **HB**: Hubei,
  **HI**: Hainan, **QH**: Qinghai, **HN**: Hunan, **JX**: Jiangxi, **SX**: Shanxi, **SN**: Shaanxi, **ZJ**: Zhejiang,
  **GS**: Gansu, **JS**: Jiangsu, **SD**: Shandong, **BJ**: Beijing, **SH**: Shanghai, **TJ**: Tianjin,
  **CQ**: Chongqing, **MO**: Macao SAR, China, **HK**: Hong Kong SAR, China, **NX**: Ningxia, **GX**: Guangxi,
  **XJ**: Xinjiang, **XZ**: Tibet, **NM**: Inner Mongolia.

* `white` - (Required, Int, NonUpdatable) Specifies the protective action.
  Valid values are as follows:
  + `0`: WAF blocks requests that hit the rule.
  + `1`: WAF allows requests that hit the rule.
  + `2`: WAF only record requests that hit the rule.
  
* `ip_type` - (Required, String, NonUpdatable) Specifies the IP type of the rule.

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the rule will be applied.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
