---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_cloud_instance"
description: |-
  Using this resource to manage a cloud WAF in HuaweiCloud.
---

# huaweicloud_waf_cloud_instance

Using this resource to manage a cloud WAF in HuaweiCloud.

## Example Usage

### Prepaid cloud WAF

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code = "professional"

  bandwidth_expack_product {
    resource_size = 1
  }
  domain_expack_product {
    resource_size = 1
  }
  rule_expack_product {
    resource_size = 1
  }

  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  auto_renew            = "true"
  enterprise_project_id = var.enterprise_project_id
}
```

### Postpaid cloud WAF (Currently only applicable to HuaweiCloud International website)

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_waf_cloud_instance" "test" {
  charging_mode         = "postPaid"
  website               = "hec-hk"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cloud WAF is located.
  If omitted, the provider-level region will be used.
  
  Changing this will create a new resource.

* `charging_mode` - (Required, String, ForceNew) Specifies the charging mode of the cloud WAF.
  The valid values are **postPaid** and **prePaid** (the yearly/monthly billing mode).

  Changing this will create a new resource.

  -> If `charging_mode` is set to **postPaid**, only `region`, `website`, `enterprise_project_id` can be specified.
  If `charging_mode` is set to **prePaid**, `website` cannot be specified.
  The postpaid charging mode is only applicable to Huawei Cloud International website.

* `website` - (Optional, String) Specifies the website to which the account belongs. Valid value is **hec-hk**.

  This parameter is required when `charging_mode` is set to **postPaid**.

* `resource_spec_code` - (Optional, String) Specifies the specification of the cloud WAF.
  The valid values are as follows:
  + **detection**: Introduction edition.
  + **professional**: Standard edition.
  + **enterprise**: Professional edition.
  + **ultimate**: Platinum edition.

  This parameter is required when `charging_mode` is set to **prePaid**.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the cloud WAF.
  Valid values are **month** and **year**.

  Changing this will create a new resource.

  This parameter is required when `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the cloud WAF.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  
  Changing this will create a new resource.

  This parameter is required when `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled.
  Valid values are **true** and **false**.

  This parameter takes effect only when `charging_mode` is set to **prePaid**.

* `bandwidth_expack_product` - (Optional, List) Specifies the configuration of the bandwidth extended packages.
  The [bandwidth_expack_product](#extended_packages) structure is documented below.

  This parameter takes effect only when `charging_mode` is set to **prePaid**.

* `domain_expack_product` - (Optional, List) Specifies the configuration of the domain extended packages.
  The [domain_expack_product](#extended_packages) structure is documented below.

  This parameter takes effect only when `charging_mode` is set to **prePaid**.

* `rule_expack_product` - (Optional, List) Specifies the configuration of the rule extended packages.
  The [rule_expack_product](#extended_packages) structure is documented below.

  This parameter takes effect only when `charging_mode` is set to **prePaid**.

-> The specification code '**detection**' does not support extended packages.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the cloud
  WAF belongs. For enterprise users, if omitted, default enterprise project will be used.

<a name="extended_packages"></a>
The `bandwidth_expack_product`, `domain_expack_product` or `rule_expack_product` block supports:

* `resource_size` - (Optional, Int) Specifies the number of extended packages.
  + For bandwidth extended packages, each package will support `1,000` QPS or `20` Mbits/s (outside HUAWEI Cloud) and
    `50` Mbits/s (inside HUAWEI Cloud) bandwidth.
  + For domain extended packages, each package will support `10` domain names (only one level-1 domain is supported).
  + For rule extended packages, each package will support `10` protection rules (only IP black/white list is supported).

  -> The `resource_size` cannot be reduced below `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cloud WAF.

* `status` - The current status of the cloud WAF.
  + `0`: Normal.
  + `1`: Frozen.
  + `2`: Deleted.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

There are two ways to import WAF cloud instance state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_cloud_instance.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_cloud_instance.test <id>/<enterprise_project_id>
```

Note that the imported state is not identical to your resource definition, due to API response reason.

For prepaid cloud WAF, the missing attributes include `enterprise_project_id`, `period_unit`, `period` and `auto_renew`.
You can ignore changes as below.

```hcl
resource "huaweicloud_waf_cloud_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      enterprise_project_id,
      period_unit,
      period,
      auto_renew,
    ]
  }
}
```

For postPaid cloud WAF, the missing attributes include `enterprise_project_id` and `website`.
You can ignore changes as below.

```hcl
resource "huaweicloud_waf_cloud_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      enterprise_project_id,
      website,
    ]
  }
}
```
