---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_cloud_instance

Using this resource to manage a cloud WAF in HuaweiCloud.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cloud WAF is located.  
  If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `resource_spec_code` - (Required, String) Specifies the specification of the cloud WAF.  
  The valid values are as follows:
  + **detection**: Introduction edition.
  + **professional**: Standard edition.
  + **enterprise**: Professional edition.
  + **ultimate**: Platinum edition.

* `charging_mode` - (Required, String, ForceNew) Specifies the charging mode of the cloud WAF.  
  The valid value is **prePaid** (the yearly/monthly billing mode).
  Changing this will create a new resource.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit of the cloud WAF.  
  Valid values are **month** and **year**.  
  Changing this will create a new resource.

* `period` - (Required, Int, ForceNew) Specifies the charging period of the cloud WAF.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.  
  Changing this will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  Valid values are **true** and **false**.

* `bandwidth_expack_product` - (Optional, List) Specifies the configuration of the bandwidth extended packages.
  The [object](#extended_packages) structure is documented below.

* `domain_expack_product` - (Optional, List) Specifies the configuration of the domain extended packages.
  The [object](#extended_packages) structure is documented below.

* `rule_expack_product` - (Optional, List) Specifies the configuration of the rule extended packages.
  The [object](#extended_packages) structure is documented below.

-> The specification code '**detection**' does not support extended packages.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the cloud
  WAF belongs.  
  Changing this will create a new resource.

<a name="extended_packages"></a>
The `bandwidth_expack_product`, `domain_expack_product` or `rule_expack_product` block supports:

* `resource_size` - (Optional, Int) Specifies the number of extended packages.
  + For bandwidth extended packages, each package will support `1,000` QPS or `20` Mbits/s (outside HUAWEI Cloud) and
    `50` Mbits/s (inside HUAWEI Cloud) bandwidth.
  + For domain extended packages, each package will support `10` domain names (only one level-1 domain is supported).
  + For rule extended packages, each package will support `10` protection rules (only IP black/white list is supported).

-> The `resource_size` cannot be reduced below `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cloud WAF.

* `status` - The current status of the cloud WAF.
  + **0**: Normal.
  + **1**: Freezen.

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
The missing attributes include `charging_mode`, `period_unit`, `period`, `auto_renew` and `enterprise_project_id`.  
You can ignore changes as below.

```hcl
resource "huaweicloud_waf_cloud_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      enterprise_project_id,
      charging_mode,
      period_unit,
      period,
      auto_renew,
    ]
  }
}
```
