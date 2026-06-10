---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_instance"
description: |-
  Manages an AAD instance resource within HuaweiCloud.
---

# huaweicloud_aad_instance

Manages an AAD instance resource within HuaweiCloud.

-> 1. The AAD instances do not support deletion and unsubscribing, please choose your purchase carefully.
  <br/>2. Deleting this resource will not delete or unsubscribe AAD instance, but will only remove the resource
  information from the tf state file.

## Example Usage

```hcl
variable "instance_name" {}

resource "huaweicloud_aad_instance" "test" {
  ip_type              = 0
  resource_region      = "north_china"
  instance_access_type = "1"
  duration             = 1
  amount               = 1
  instance_name        = var.instance_name
  period_type          = 2
  service_bandwidth    = 100
  basic_bandwidth      = 10
  elastic_bandwidth    = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `ip_type` - (Required, Int, NonUpdatable) Specifies the IP type.  
  The valid values are as follows:
  + **0**: IPv4.
  + **1**: IPv6.

  When `resource_region` is **asia_pacific**, only **0** is supported.

* `resource_region` - (Required, String, NonUpdatable) Specifies the resource region.  
  The valid values are as follows:
  + **north_china**
  + **east_china**
  + **asia_pacific**

* `instance_access_type` - (Required, String, NonUpdatable) Specifies the instance access type.  
  The valid values are as follows:
  + **0**: Website.
  + **1**: IP access.

* `duration` - (Required, Int, NonUpdatable) Specifies the subscription duration.
  The valid values are **1** to **9** for monthly subscription and **1** for yearly subscription.

* `amount` - (Required, Int, NonUpdatable) Specifies the number of instances to purchase.

* `instance_name` - (Required, String, NonUpdatable) Specifies the instance name.

* `period_type` - (Required, Int, NonUpdatable) Specifies the subscription period type.  
  The valid values are as follows:
  + **2**: Month.
  + **3**: Year.

* `service_bandwidth` - (Required, Int) Specifies the service bandwidth.

* `basic_bandwidth` - (Optional, Int) Specifies the basic bandwidth.  
  + When `resource_region` is **north_china** or **east_china**, `basic_bandwidth` cannot be empty.
  + When `resource_region` is **asia_pacific**, `basic_bandwidth` must be empty.

* `elastic_bandwidth` - (Optional, Int) Specifies the elastic bandwidth.  
  + When `resource_region` is **north_china** or **east_china**, `elastic_bandwidth` cannot be empty.
  + When `resource_region` is **asia_pacific**, `elastic_bandwidth` must be empty.

* `basic_qps` - (Optional, Int) Specifies the service QPS.  
  + When `instance_access_type` is `1`, `basic_qps` must be empty.
  + When `instance_access_type` is `0`, `basic_qps` cannot be empty.
  + When `resource_region` is **north_china** or **east_china**, the value range of `basic_qps` is `3000` to `100000`.
  + When `resource_region` is **asia_pacific**, the value range of `basic_qps` is `1000` to `100000`.

* `elastic_service_bandwidth_type` - (Optional, Int) Specifies the elastic service bandwidth type, empty means not
  turned on.  
  The valid values are as follows:
  + **2**: Daily 95th percentile.
  + **3**: Monthly 95th percentile.

  + When `resource_region` is **north_china** or **east_china**, only **3** is supported.

* `elastic_service_bandwidth` - (Optional, Int) Specifies the elastic service bandwidth increment.

* `protection_package` - (Optional, String, NonUpdatable) Specifies the protection package.  
  The valid values are as follows:
  + **basic**: Insurance protection.
  + **unlimited**: Unlimited protection.

  + When `resource_region` is **asia_pacific**, `protection_package` cannot be empty.
  + When `resource_region` is not **asia_pacific**, `protection_package` must be empty.

* `protected_domain` - (Optional, Int) Specifies the number of protected domains.  
  + When `instance_access_type` is `1`, `protected_domain` must be empty.
  + When `instance_access_type` is `0`, `protected_domain` cannot be empty.
  + When `resource_region` is **north_china** or **east_china**, the value range of `protected_domain` is `50` to `500`.
  + When `resource_region` is **asia_pacific**, the value range of `protected_domain` is `10` to `200`.

* `forwarding_rule` - (Optional, Int) Specifies the number of forwarding rules.  
  + When `instance_access_type` is `1`, `forwarding_rule` cannot be empty.
  + When `instance_access_type` is `0`, `forwarding_rule` must be empty.
  + When `resource_region` is **north_china** or **east_china**, the value range of `forwarding_rule` is `50` to `500`.
  + When `resource_region` is **asia_pacific**, the value range of `forwarding_rule` is `5` to `200`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `port_num` - (Optional, Int) Specifies the number of ports.

* `bind_domain_num` - (Optional, Int) Specifies the number of bound domains.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The instance ID.

* `version` - The version.

* `expire_time` - The expiration time.

* `create_time` - The creation time.

* `current_time` - The current time.

* `product_uuid` - The product specification UUID.

* `isp_spec` - The ISP specification.

* `data_center` - The data center.

* `spec_type` - The product specification type.

* `main_resource_type` - The main resource type.

* `main_resource_spec_code` - The main resource specification code.

* `main_resource_product_id` - The main resource product ID.

* `instance_config` - The instance configuration.

  The [instance_config](#instance_config_struct) structure is documented below.

* `elastic_service_bw_update_enable` - Whether the elastic service bandwidth can be updated.

<a name="instance_config_struct"></a>
The `instance_config` block supports:

* `tags` - The instance-level tags.

* `freeze_type` - The freeze type list.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Defaults to `30` minutes.
* `update` - Defaults to `30` minutes.

## Import

The AAD instances can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aad_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `ip_type`, `resource_region`, `instance_access_type`, `duration`,
`amount`, `period_type`, `basic_qps`, `protection_package`, `protected_domain`, and `forwarding_rule`.
It is generally recommended running `terraform plan` after importing a resource. You can keep the resource the same with
its definition by choosing any of them to update. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_aad_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      ip_type,
      resource_region,
      instance_access_type,
      duration,
      amount,
      period_type,
      basic_qps,
      protection_package,
      protected_domain,
      forwarding_rule,
    ]
  }
}
```
