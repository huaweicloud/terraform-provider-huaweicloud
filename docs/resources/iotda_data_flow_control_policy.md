---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_data_flow_control_policy"
description: |-
  Manages an IoTDA data flow control policy resource within HuaweiCloud.
---

# huaweicloud_iotda_data_flow_control_policy

Manages an IoTDA data flow control policy resource within HuaweiCloud.

-> Currently, data flow control policy resources are only supported on IoTDA **standard** or **enterprise** edition
  instance. When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify
  the IoTDA service endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  *9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com*, then you need to configure the
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
variable "name" {}
variable "scope" {}

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name  = var.name
  scope = var.scope
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Optional, String) Specifies the name of the data flow control policy. The length must not exceed `256`, and
  only Chinese characters, letters, numbers, and the following characters are allowed: `_?'#().,&%@!-**`.

* `description` - (Optional, String) Specifies the description of the data flow control policy. The length must not
  exceed `256`, and only Chinese characters, letters, numbers, and the following characters are allowed: `_?'#().,&%@!-**`.

* `scope` - (Optional, String, ForceNew) Specifies the scope of the data flow control policy. Changing this parameter
  will create a new resource.  
  The valid values are as follows:
  + **USER**: Tenant level flow control strategy.
  + **CHANNEL**: Forwarding channel level flow control strategy.
  + **RULE**: Forwarding rule level flow control strategy.
  + **ACTION**: Forwarding action level flow control strategy.

  If omitted, defaults to **USER**.

* `scope_value` - (Optional, String, ForceNew) Specifies the scope add value of the data flow control policy. Changing
  this parameter will create a new resource.  
  + If the `scope` is set to **USER**, this field does not need to be set.
  + If the `scope` is set to **CHANNEL**, the valid values are **HTTP_FORWARDING**, **DIS_FORWARDING**,
    **OBS_FORWARDING**, **AMQP_FORWARDING**, and **DMS_KAFKA_FORWARDING**.
  + If the `scope` is set to **RULE**, the value of this field is the corresponding data forwarding rule ID.
  + If the `scope` is set to **ACTION**, the value of this field is the corresponding data forwarding rule action ID.

* `limit` - (Optional, Int) Specifies the size of the data forwarding flow control, in tps. Integers with valid values
  ranging from `1` to `1,000`. Defaults to `1,000`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The data flow control policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_data_flow_control_policy.test <id>
```
