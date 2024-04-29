---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_dataforwarding_rules"
description: ""
---

# huaweicloud_iotda_dataforwarding_rules

Use this data source to get the list of IoTDA dataforwarding rules.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
  endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  **9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com**, then you need to configure the
  `provider` block as follows:

  ```hcl
  provider "huaweicloud" {
    endpoints = {
      iotda = "https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com"
    }
  }
  ```

## Example Usage

```hcl
variable "rule_name" {}

data "huaweicloud_iotda_dataforwarding_rules" "test" {
  name = var.rule_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the dataforwarding rules.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the dataforwarding rule.

* `name` - (Optional, String) Specifies the name of the dataforwarding rule.

* `resource` - (Optional, String) Specifies the data source of the dataforwarding rule.
  This parameter must be used together with `trigger`. The valid values are as follows:
  + **device**
  + **device.property**
  + **device.message**
  + **device.message.status**
  + **device.status**
  + **batchtask**
  + **product**
  + **device.command.status**

* `trigger` - (Optional, String) Specifies the triggering event of the data source corresponding to
  the dataforwarding rule. This parameter must be used together with `resource`. The valid values are as follows:
  + **device:create**: Device added.
  + **device:delete**: Device deleted.
  + **device:update**: Device updated.
  + **device.status:update**: Device status changed.
  + **device.property:report**: Device property reported.
  + **device.message:report**: Device message reported.
  + **device.message.status:update**: Device message status changed.
  + **batchtask:update**: Batch task status changed.
  + **product:create**: Product added.
  + **product:delete**: Product deleted.
  + **product:update**: Product updated.
  + **device.command.status:update**: Device asynchronous command status updated.

* `app_type` - (Optional, String) Specifies the validity scope of the dataforwarding rule.
  The valid values are as follows:
  + **GLOBAL**: The validity scope is tenant level.
  + **APP**: The validity scope is resource space level.

  -> If the `app_type` value is **APP**, this parameter can be used together with the `space_id` to query
    the dataforwarding rules in the corresponding resource space, if not associated with the `space_id`,
    will be query the dataforwarding rules in the default resource space.

* `space_id` - (Optional, String) Specifies the ID of the resource space to which the dataforwarding rule belongs.

  -> If use this parameter to query, the parameter `app_type` must be set to **APP**.

* `enabled` - (Optional, String) Specifies whether to enable the dataforwarding rule.
  The value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - All rules that match the filter parameters.
  The [rules](#iotda_rules) structure is documented below.

<a name="iotda_rules"></a>
The `rules` block supports:

* `id` - The ID of the dataforwarding rule.

* `name` - The name of the dataforwarding rule.

* `description` - The description of the dataforwarding rule.

* `resource` - The data source of the dataforwarding rule.

* `trigger` - The triggering event of the data source corresponding to the dataforwarding rule.

* `app_type` - The validity scope of the dataforwarding rule.

* `space_id` - The ID of the resource space to which the dataforwarding rule belongs.

* `enabled` - Whether to enable the dataforwarding rule.

* `select` - The user defined SQL **select** statement.

* `where` - The user defined SQL **where** statement.
