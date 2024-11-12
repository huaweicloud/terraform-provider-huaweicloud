---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_geolocation_access_control"
description: |-
  Manages a WAF rule geolocation access control resource within HuaweiCloud.
---

# huaweicloud_waf_rule_geolocation_access_control

Manages a WAF rule geolocation access control resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The geolocation access control rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable policy_id {}
variable enterprise_project_id {}

resource "huaweicloud_waf_rule_geolocation_access_control" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  name                  = "test_rule"
  geolocation           = "FJ|JL|LN|GZ"
  action                = 1
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of WAF geolocation access control rule. This parameter can contain a
  maximum of `128` characters. Only letters, digits, hyphens (-), underscores (_), colons (:) and periods (.) are allowed.

* `geolocation` - (Required, String) Specifies the locations that can be configured in the geolocation access control
  rule. Separate multiple locations with vertical lines, such as **FJ|JL|LN|GZ**.

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

* `action` - (Required, Int) Specifies the protective action of WAF geolocation access control rule.
  Valid values are as follows:
  + `0`: WAF blocks requests that hit the rule.
  + `1`: WAF allows requests that hit the rule.
  + `2`: WAF only record requests that hit the rule.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF geolocation access
  control rule. For enterprise users, if omitted, default enterprise project will be used.

  Changing this parameter will create a new resource.

* `status` - (Optional, Int) Specifies the status of WAF geolocation access control rule.
  Valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

  The default value is `1`.

* `description` - (Optional, String) Specifies the description of WAF geolocation access control rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

There are two ways to import WAF rule geolocation access control state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_geolocation_access_control.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_geolocation_access_control.test <policy_id>/<rule_id>/<enterprise_project_id>
```
