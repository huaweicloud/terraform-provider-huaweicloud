---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_group_resource_relations"
description: |-
  Use this data source to get the list of COC relationships between groups and resources.
---

# huaweicloud_coc_group_resource_relations

Use this data source to get the list of COC relationships between groups and resources.

## Example Usage

```hcl
variable "application_id" {}

data "huaweicloud_coc_group_resource_relations" "test" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `cloud_service_name` - (Required, String) Specifies the cloud service name.
  The value can be **ecs**, **cce**, **rds** and so on.

* `type` - (Required, String) Specifies the resource type name.

* `vendor` - (Required, String) Specifies the manufacturer information.
  Values can be as follows:
  + **RMS**: Huawei Cloud Vendor.
  + **ALI**: Alibaba Cloud Vendor.
  + **OTHER**: Other Vendor.

* `application_id` - (Optional, String) Specifies the application ID associated with the group.

* `component_id` - (Optional, String) Specifies the component ID associated with the group.

* `group_id` - (Optional, String) Specifies the group ID.

-> `application_id`, `component_id` and `group_id` cannot coexist, and one of them is required.

* `resource_id_list` - (Optional, List) Specifies the resource ID list.

* `name` - (Optional, String) Specifies the cloud resource name.

* `region_id` - (Optional, String) Specifies the region ID.

* `az_id` - (Optional, String) Specifies the availability zone ID.

* `ip_type` - (Optional, String) Specifies the IP type.
  Values can be as follows:
  + **fixed**: Intranet IP.
  + **floating**: Elastic public IP.

* `ip` - (Optional, String) Specifies the cloud resource IP.

* `status` - (Optional, String) Specifies the resource status.
  For details, see [status](https://support.huaweicloud.com/api-ecs/ecs_08_0002.html).

* `agent_state` - (Optional, String) Specifies the unified agent status.
  Values can be **ONLINE**, **OFFLINE**, **INSTALLING**, **FAILED**, **UNINSTALLED** or **null**.

* `image_name` - (Optional, String) Specifies the fuzzy query the image name.

* `os_type` - (Optional, String) Specifies the cloud resource operating system type.
  Values can be **windows** or **linux**.

* `tag` - (Optional, String) Specifies the tags for cloud resources. The format of the tag is **key.value**.
  When naming tags, the following requirements must be met:
  + The **key** of the tag can only contain uppercase letters (A~Z), lowercase letters (a~z), numbers (0-9),
    underscores (\_), hyphens (-), and Chinese characters.
  + The **value** of the tag can only contain uppercase letters (A~Z),lowercase letters (a~z), numbers (0-9),
    underscores (\_), hyphens (-), decimal points (.), and Chinese characters.

* `charging_mode` - (Optional, String) Specifies the billing type for the cloud server.
  Values can be as follows:
  + **0**: On demand billing.
  + **1**: Yearly package or monthly package.
  + **2**: Award based billing.

* `flavor_name` - (Optional, String) Specifies the cloud resource specification name.

* `ip_list` - (Optional, List) Specifies the cloud resource IP list.

* `is_collected` - (Optional, Bool) Specifies whether it is a favorite.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of relationships between groups and resources.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the application ID.

* `cmdb_resource_id` - Indicates the cmdb database resource ID.

* `group_id` - Indicates the group ID.

* `group_name` - Indicates the group name.

* `resource_id` - Indicates the resource ID.

* `name` - Indicates the resource name.

* `cloud_service_name` - Indicates the cloud service name.

* `type` - Indicates the resource type.

* `region_id` - Indicates the region ID.

* `ep_id` - Indicates the enterprise project ID.

* `ep_name` - Indicates the enterprise project name.

* `project_id` - Indicates the project ID.

* `domain_id` - Indicates the tenant ID.

* `tags` - Indicates the resource tags.

  The [tags](#data_tags_struct) structure is documented below.

* `agent_id` - Indicates the ID assigned by unified agent.

* `agent_state` - Indicates the unified agent status.

* `inner_ip` - Indicates the inner IP.

* `properties` - Indicates the resource properties.

* `ingest_properties` - Indicates the ingest attributes of the resource.

* `operable` - Indicates whether the user defined resource can operate the instance.

* `create_time` - Indicates the creation time.

<a name="data_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.
