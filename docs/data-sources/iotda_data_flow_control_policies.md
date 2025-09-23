---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_data_flow_control_policies"
description: |-
  Use this data source to get the list of IoTDA data flow control policies within HuaweiCloud.
---

# huaweicloud_iotda_data_flow_control_policies

Use this data source to get the list of IoTDA data flow control policies within HuaweiCloud.

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
variable "policy_name" {}

data "huaweicloud_iotda_data_flow_control_policies" "test" {
  policy_name = var.policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data flow control policies.
  If omitted, the provider-level region will be used.

* `scope` - (Optional, String) Specifies the scope of the data flow control policies.  
  The valid values are as follows:
  + **USER**: Tenant level flow control strategy.
  + **CHANNEL**: Forwarding channel level flow control strategy.
  + **RULE**: Forwarding rule level flow control strategy.
  + **ACTION**: Forwarding action level flow control strategy.

  If omitted, query all scope data flow control policies.

* `scope_value` - (Optional, String) Specifies the scope add value of the data flow control policies.
  + If omitted or the `scope` is set to **USER**, this field does not need to be set, representing the query of tenant
    level flow control policies.
  + If the `scope` is set to **CHANNEL**, the valid values are **HTTP_FORWARDING**, **DIS_FORWARDING**,
    **OBS_FORWARDING**, **AMQP_FORWARDING**, and **DMS_KAFKA_FORWARDING**. If omitted, query all forwarding channel
    level flow control policies.
  + If the `scope` is set to **RULE**, the value of this field is the corresponding rule ID. If omitted, query all
    forwarding rule level flow control policies.
  + If the `scope` is set to **ACTION**, the value of this field is the corresponding rule action ID. If omitted, query
    all forwarding action level flow control policies.

  -> The `scope_value` must be used together with `scope` and is invalid when used alone.

* `policy_name` - (Optional, String) Specifies the name of the data flow control policy. This field will undergo a fuzzy
  matching query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `policies` - All data flow control policies that match the filter parameters.
  The [policies](#iotda_policies) structure is documented below.

<a name="iotda_policies"></a>
The `policies` block supports:

* `id` - The ID of the data flow control policy.

* `name` - The name of the data flow control policy.

* `description` - The description of the data flow control policy.

* `scope` - The scope of the data flow control policy.

* `scope_value` - The scope add value of the data flow control policy.

* `limit` - The size of the data forwarding flow control, in tps.
