---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instances"
description: |-
  Use this data source to query the instances within HuaweiCloud.
---

# huaweicloud_apig_instances

Use this data source to query the instances within HuaweiCloud.

## Example Usage

```hcl
variable "instance_name" {}

data "huaweicloud_apig_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance to be queried.

* `name` - (Optional, String) Specifies the name of the instance to be queried.

* `status` - (Optional, String) Specifies the status of the instance to be queried.  
  The valid values are as follows:
  + **Creating**: Instance creation in progress.
  + **CreateSuccess**: Instance created successfully
  + **CreateFail**: Instance creation failed.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the instances belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - All instances that match the filter parameters.
  The [instances](#attrblock_instances) structure is documented below.

<a name="attrblock_instances"></a>
The `instances` block supports:

* `id` - The ID of instance.

* `name` - The name of instance.

* `type` - The type of instance.

* `status` - The status of instance.

* `edition` - The edition of instance.

* `eip_address` - The elastic IP address of instance binding.

* `enterprise_project_id` - The enterprise project ID of the instance.

* `loadbalancer_provider` - The type of load balancer used by the instance.  
  The valid values are as follows:
  + **lvs**: Linux virtual server.
  + **elb**: Elastic load balance.

* `created_at` - The creation time of the instance, in RFC3339 format.
