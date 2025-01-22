---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_amqps"
description: |-
  Use this data source to get the list of the IoTDA AMQP queues within HuaweiCloud.
---

# huaweicloud_iotda_amqps

Use this data source to get the list of the IoTDA AMQP queues within HuaweiCloud.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify
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
variable "queue_name" {}

data "huaweicloud_iotda_amqps" "test" {
  name = var.queue_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the AMQP queues.
  If omitted, the provider-level region will be used.

* `queue_id` - (Optional, String) Specifies the ID of the AMQP queue.

* `name` - (Optional, String) Specifies the name of the AMQP queue.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `queues` - The list of the AMQP queues.
  The [queues](#iotda_queues) structure is documented below.

<a name="iotda_queues"></a>
The `queues` block supports:

* `id` - The ID of the AMQP queue.

* `name` - The name of the AMQP queue.

* `created_at` - The creation time of the AMQP queue.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

* `updated_at` - The latest update time of the AMQP queue.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.
