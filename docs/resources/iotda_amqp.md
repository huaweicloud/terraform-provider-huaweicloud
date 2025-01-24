---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_amqp"
description: |-
  Manages an IoTDA AMQP queue resource within HuaweiCloud.
---

# huaweicloud_iotda_amqp

Manages an IoTDA AMQP queue resource within HuaweiCloud.

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
variable "queue_name" {}

resource "huaweicloud_iotda_amqp" "test" {
  name = var.queue_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA AMQP queue resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the AMQP queue name, which contains `8` to `128` characters.
  Only letters, digits, hyphens (-), underscores (_), dots (.) and colons (:) are allowed.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The creation time of the AMQP queue.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

* `updated_at` - The latest update time of the AMQP queue.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

## Import

The AMQP queue can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_amqp.test <id>
```
