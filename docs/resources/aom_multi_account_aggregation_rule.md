---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_multi_account_aggregation_rule"
description:  |-
  Manages an AOM multi account aggregation rule resource within HuaweiCloud.
---

# huaweicloud_aom_multi_account_aggregation_rule

Manages an AOM multi account aggregation rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "account_id" {}
variable "account_name" {}

resource "huaweicloud_aom_multi_account_aggregation_rule" "test" {
  instance_id = var.instance_id

  accounts {
    id   = var.account_id
    name = var.account_name
  }

  services {
    service = "SYS.ELB"
    metrics = [
        "huaweicloud_sys_elb_m1_cps",
        "huaweicloud_sys_elb_m2_act_conn",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the AOM prometheus instance ID.
  Changing this parameter will create a new resource.

* `accounts` - (Required, List) Specifies the accounts list.
  The [accounts](#accounts_struct) structure is documented below.

* `services` - (Optional, List) Specifies the services list.
  The [services](#services_struct) structure is documented below.

* `send_to_source_account` - (Optional, Bool) Specifies whether the member accounts retain metric data after they are
  connected to the prometheus instance for aggregation.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `id` - (Required, String) Specifies the account ID.

* `name` - (Required, String) Specifies the account name.

* `urn` - (Optional, String) Specifies the uniform resource name of the account.

* `join_method` - (Optional, String) Specifies the method how the account joined in the organization.

* `joined_at` - (Optional, String) Specifies the time when the account joined in the organization.

<a name="services_struct"></a>
The `services` block supports:

* `service` - (Required, String) Specifies the service name.

* `metrics` - (Required, List) Specifies the metrics List.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID which is same as `instance_id`.

## Import

The AOM multi account aggregation rule resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_multi_account_aggregation_rule.test <id>
```
