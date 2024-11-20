---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_data_backlog_policy"
description: |-
  Manages an IoTDA data backlog policy resource within HuaweiCloud.
---

# huaweicloud_iotda_data_backlog_policy

Manages an IoTDA data backlog policy resource within HuaweiCloud.

-> 1. A single tenant can create up to one data backlog policy under a single IoTDA instance.
  <br/>2. Before creating a data backlog policy, it is necessary to ensure that there is a data forwarding rule under
  the current IoTDA instance.
  <br/>3. After the successful creation of the data backlog policy, it will take effect for all data forwarding rules,
  covering the default backlog size (`1` GB) and default backlog time (`1` day) of all forwarding rules.
  <br/>4. When the maximum backlog (cache) size or backlog (cache) time is exceeded, the earliest unprocessed flow data
  will be discarded until the backlog (cache) size and time limits are met. Please consider carefully before using a
  data backlog policy and configuring a reasonable backlog size and time.

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
variable "name" {}

resource "huaweicloud_iotda_data_backlog_policy" "test" {
  name  = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Optional, String) Specifies the name of the data backlog policy. The length must not exceed `256`, and
  only Chinese characters, letters, numbers, and the following characters are allowed: `_?'#().,&%@!-`.

* `description` - (Optional, String) Specifies the description of the data backlog policy. The length must not exceed
  `256`, and only Chinese characters, letters, numbers, and the following characters are allowed: `_?'#().,&%@!-`.

* `backlog_size` - (Optional, String) Specifies the size of data backlog in bytes. The range of valid values is integers
  from `0` to `1,073,741,823`, defaults to `1,073,741,823` (`1` GB).
  + When `backlog_size` is set to `0`, it means there is no backlog.

* `backlog_time` - (Optional, String) Specifies the data backlog time in seconds. The range of valid values is integers
  from `0` to `86,399`, defaults to `86,399` (`1` day).
  + When `backlog_time` is set to `0`, it means there is no backlog.

-> If both `backlog_size` and `backlog_time` dimensions are configured, the dimension that reaches the threshold first
   shall prevail.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The data backlog policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_data_backlog_policy.test <id>
```
