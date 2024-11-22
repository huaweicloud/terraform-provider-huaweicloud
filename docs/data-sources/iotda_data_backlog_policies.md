---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_data_backlog_policies"
description: |-
  Use this data source to get the list of IoTDA data backlog policies within HuaweiCloud.
---

# huaweicloud_iotda_data_backlog_policies

Use this data source to get the list of IoTDA data backlog policies within HuaweiCloud.

-> Currently, data backlog policy resources are only supported on IoTDA **standard** or **enterprise** edition
  instance. When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify
  the IoTDA service endpoint in `provider` block.
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
variable "policy_name" {}

data "huaweicloud_iotda_data_backlog_policies" "test" {
  policy_name = var.policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data backlog policies.
  If omitted, the provider-level region will be used.

* `policy_name` - (Optional, String) Specifies the name of the data backlog policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `policies` - All data backlog policies that match the filter parameters.
  The [policies](#iotda_policies) structure is documented below.

<a name="iotda_policies"></a>
The `policies` block supports:

* `id` - The ID of the data backlog policy.

* `name` - The name of the data backlog policy.

* `description` - The description of the data backlog policy.

* `backlog_size` - The size of data backlog in bytes. The range of valid values is integers from `0` to
  `1,073,741,823` (`1` GB), `0` means no backlog.

* `backlog_time` - The data backlog time in seconds. The range of valid values is integers from `0` to
  `86,399` (`1` day), `0` means no backlog.
