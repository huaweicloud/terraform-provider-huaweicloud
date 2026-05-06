---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_devserver_flavors"
description: |-
  Use this data source to get the list of available DevServer flavors of ModelArts within Huaweicloud.
---

# huaweicloud_modelarts_devserver_flavors

Use this data source to get the list of available DevServer flavors of ModelArts within Huaweicloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_modelarts_devserver_flavors" "test" {}
```

### Filter by server type

```hcl
data "huaweicloud_modelarts_devserver_flavors" "test" {
  server_type = "ECS"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DevServer flavors are located.  
  If omitted, the provider-level region will be used.

* `server_type` - (Optional, String) Specifies the service type of the DevServer flavors.  
  The valid values are as follows:
  + **BMS**: Bare metal server.
  + **ECS**: Elastic cloud server.
  + **HPS**: Hyper node server.

* `arch` - (Optional, String) Specifies the CPU architecture of the DevServer flavors.
  The valid values are as follows:
  + **X86**
  + **ARM**

* `charging_mode` - (Optional, String) Specifies the charging mode of the DevServer flavors.  
  The valid values are as follows:
  + **PRE_PAID**: The prepaid mode.
  + **POST_PAID**: The postpaid mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of DevServer flavors that matched filter parameters.  
  The [flavors](#devserver_flavors_attr) structure is documented below.

<a name="devserver_flavors_attr"></a>
The `flavors` block supports:

* `flavor` - The IaaS flavor name.

* `specification` - The specification configuration of the flavor.

* `arch` - The CPU architecture of the flavor.

* `server_type` - The server type of the flavor.

* `sku_code` - The SKU billing code of the flavor.

* `charging_mode` - The charging mode of the flavor.

* `roce_num` - The number of NICs of the flavor.

* `count` - The number of super node instances.

* `status` - The status of the server flavor.

* `server_flavors` - The expandable flavors of the super node.  
  The [server_flavors](#devserver_flavors_server_flavors) structure is documented below.

* `flavor_type` - The computing power card type of the flavor.

* `availability_zones` - The availability zones of the flavor.  
  The [availability_zones](#devserver_flavors_availability_zones) structure is documented below.

<a name="devserver_flavors_server_flavors"></a>
The `server_flavors` block supports:

* `name` - The name of the server flavor.

<a name="devserver_flavors_availability_zones"></a>
The `availability_zones` block supports:

* `id` - The ID of the availability zone.

* `is_sold_out` - Whether the flavor is sold out in the availability zone.
