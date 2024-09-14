---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_dataforwarding_rule"
description: ""
---

# huaweicloud_iotda_dataforwarding_rule

Manages an IoTDA data forwarding rule within HuaweiCloud.

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
variable "disRegion" {}
variable "disStreamId" {}
variable "obsRegion" {}
variable "obsbucket" {}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "ruleName"
  trigger = "product:create"
  enabled = true

  targets {
    type = "DIS_FORWARDING"
    dis_forwarding {
      region    = var.disRegion
      stream_id = var.disStreamId
    }
  }

  targets {
    type = "HTTP_FORWARDING"
    http_forwarding {
      url = "http://www.yourDomain.com"
    }
  }

  targets {
    type = "OBS_FORWARDING"
    obs_forwarding {
      region = var.obsRegion
      bucket = var.obsbucket
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA data forwarding rule
  resource. If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of data forwarding rule. The name contains a maximum of `256` characters.
  Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special characters are
  allowed: `?'#().,&%@!`.

* `trigger` - (Required, String, ForceNew) Specifies the trigger event. The options are as follows:
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
  + **device.command.status:update**: Update of the device asynchronous command status.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of data forwarding rule. The description contains
  a maximum of `256` characters.

* `space_id` - (Optional, String, ForceNew) Specifies the resource space ID which uses the data forwarding rule.
  If omitted, all resource space will use the data forwarding rule. Changing this parameter will create a new resource.

* `select` - (Optional, String) Specifies the SQL SELECT statement which contains a maximum of `500` characters.

* `where` - (Optional, String) Specifies the SQL WHERE statement which contains a maximum of `500` characters.

* `targets` - (Optional, List) Specifies the list of the targets (HUAWEI CLOUD services or private servers) to which you
  want to forward the data. The [targets](#IoTDA_targets) structure is documented below.

-> The IoTDA supports connection with other cloud services of HUAWEI CLOUD. When creating a rule for connecting to
  DIS, OBS, Kafka, ROMA Connect service, and SMN for the first time, you need to create a agency which can access the
  target cloud services. The default agency is `iotda_admin_trust`.

* `enabled` - (Optional, Bool) Specifies whether to enable the data forwarding rule. Defaults to **false**.
  Can not enable without `targets`.

<a name="IoTDA_targets"></a>
The `targets` block supports:

* `type` - (Required, String) Specifies the type of forwarding target. The options are as follows:
  + **HTTP_FORWARDING**: The platform can push specified device data to a Third-party application (HTTP push).
   You can set different addresses that different types of device data are pushed to.
  + **DIS_FORWARDING**: DIS provides efficient collection, transmission, and distribution of real-time data. It also
   provides an abundant selection of APIs to help you quickly create real-time data applications.
  + **OBS_FORWARDING**: OBS is a stable, secure, cloud storage service that is scalable, efficient and easy-to-use.
   It allows you to store any amount of unstructured data in any format, and provides REST APIs so you can access your
   data from anywhere.
  + **AMQP_FORWARDING**: AMQP provides a scalable, distributed message queue that supports high throughput with low
   latency. AMQP is ready from the get-go and is O&M free.
  + **DMS_KAFKA_FORWARDING**: Distributed Message Service (DMS) for Kafka features high throughput, concurrency, and
   scalability. It is suitable for real-time data transmission, stream data processing, system decoupling,
   and traffic balancing.
  + **FUNCTIONGRAPH_FORWARDING**: By forwarding data to FunctionGraph service, you only need to write your business
    function code and set the conditions for execution in FunctionGraph. There is no need to configure and manage
    servers or other infrastructure. Functions will run in an elastic, maintenance-free, and highly reliable manner.
    Currently, only standard and enterprise edition IoTDA instances are supported.

* `http_forwarding` - (Optional, List) Specifies the detail of the HTTP forwards. It is required when type
  is `HTTP_FORWARDING`. The [http_forwarding](#IoTDA_http_forwarding) structure is documented below.

* `dis_forwarding` - (Optional, List) Specifies the detail of the DIS forwards. It is required when type
  is `DIS_FORWARDING`. The [dis_forwarding](#IoTDA_dis_forwarding) structure is documented below.

* `obs_forwarding` - (Optional, List) Specifies the detail of the OBS forwards. It is required when type
  is `OBS_FORWARDING`. The [obs_forwarding](#IoTDA_obs_forwarding) structure is documented below.

* `amqp_forwarding` - (Optional, List) Specifies the detail of AMQP forwards. It is required when type
  is `AMQP_FORWARDING`. The [amqp_forwarding](#IoTDA_amqp_forwarding) structure is documented below.

* `kafka_forwarding` - (Optional, List) Specifies the detail of the KAFKA forwards. It is required when type
  is `DMS_KAFKA_FORWARDING`. The [properties](#IoTDA_kafka_forwarding) structure is documented below.

* `fgs_forwarding` - (Optional, List) Specifies the detail of the FunctionGraph forwards. It is required when
  type is **FUNCTIONGRAPH_FORWARDING**. The [fgs_forwarding](#IoTDA_fgs_forwarding) structure is documented below.

<a name="IoTDA_http_forwarding"></a>
The `http_forwarding` block supports:

* `url` - (Required, String) Specifies the Push URL. The request method must is post.

<a name="IoTDA_dis_forwarding"></a>
The `dis_forwarding` block supports:

* `region` - (Required, String) Specifies the region to which the DIS stream belongs.

* `stream_id` - (Required, String) Specifies the DIS stream ID.

* `project_id` - (Optional, String) Specifies the project ID to which the DIS stream belongs.
If omitted, the default project in the region will be used.

<a name="IoTDA_obs_forwarding"></a>
The `obs_forwarding` block supports:

* `region` - (Required, String) Specifies the region to which the OBS belongs.

* `bucket` - (Required, String) Specifies the OBS Bucket.

* `project_id` - (Optional, String) Specifies the project ID to which the OBS belongs.
  If omitted, the default project in the region will be used.

* `custom_directory` - (Optional, String) Specifies the custom directory for storing channel files. The ID contains a
 maximum of `256` characters. Multi-level directories can be separated by (/), and cannot start or end with a slash (/),
 and cannot contain more than two adjacent slashes (/). Only letters, digits, hyphens (-), underscores (_), slash (/)
 and braces ({}) are allowed. Braces can be used only for the time template parameters. For example, if the custom
 directory is in the format of {YYYY}/{MM}/{DD}/{HH}, data is generated in the directory based on the current
 time(for example, 2022/06/14/10) when data is forwarded.

<a name="IoTDA_amqp_forwarding"></a>
The `amqp_forwarding` block supports:

* `queue_name` - (Required, String) Specifies the AMQP Queue name.

<a name="IoTDA_kafka_forwarding"></a>
The `kafka_forwarding` block supports:

* `region` - (Required, String) Specifies the region to which the KAFKA belongs.

* `topic` - (Required, String) Specifies the topic.

* `addresses` - (Required, List) Specifies the list of the connected service addresses.
The [addresses](#IoTDA_forwarding_addresses) structure is documented below.

* `project_id` - (Optional, String) Specifies the project ID to which the KAFKA belongs.
If omitted, the default project in the region will be used.

* `user_name` - (Optional, String) Specifies the SASL user name.

* `password` - (Optional, String) Specifies the password.

<a name="IoTDA_fgs_forwarding"></a>
The `fgs_forwarding` block supports:

* `func_urn` - (Required, String) Specifies the function URN.

* `func_name` - (Required, String) Specifies the function name.

<a name="IoTDA_forwarding_addresses"></a>
The `addresses` block supports:

* `port` - (Required, Int) Specifies the port of the connected service address.

* `ip` - (Optional, String) Specifies the IP of the connected service address.
Exactly one of `ip` or `domain` must be provided.

* `domain` - (Optional, String) Specifies the domain of the connected service address.
Exactly one of `ip` or `domain` must be provided.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `targets/id` - The ID of the data forwarding target.

## Import

Data forwarding rules can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_dataforwarding_rule.test 10022532f4f94f26b01daa1e424853e1
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password` of `kafka_forwarding`. It is
generally recommended running `terraform plan` after importing the resource. You can then decide if changes should be
applied to the resource, or the resource definition should be updated to align with the group. Also you can ignore
changes as below.

```hcl
resource "huaweicloud_iotda_device_group" "test" {
    ...

  lifecycle {
    ignore_changes = [
      targets,
    ]
  }
}
```
