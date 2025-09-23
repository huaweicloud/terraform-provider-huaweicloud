---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_resources"
description: |-
  Use this data source to get the list of COC resources.
---

# huaweicloud_coc_resources

Use this data source to get the list of COC resources.

## Example Usage

```hcl
variable "cloud_service_name" {}
variable "type" {}

data "huaweicloud_coc_resources" "test" {
  cloud_service_name = var.cloud_service_name
  type               = var.type
}
```

## Argument Reference

The following arguments are supported:

* `cloud_service_name` - (Required, String) Specifies the cloud service name.

* `type` - (Required, String) Specifies the resource type name.

* `name` - (Optional, String) Specifies the cloud resource name.

* `ep_id` - (Optional, String) Specifies the enterprise project ID.

* `project_id` - (Optional, String) Specifies the project ID.

* `region_id` - (Optional, String) Specifies the region ID.

* `az_id` - (Optional, String) Specifies the availability zone ID.

* `ip_type` - (Optional, String) Specifies the IP type.
  Values can be as follows:
  + **fixed**: Intranet IP.
  + **floating**: Elastic public IP.

* `ip` - (Optional, String) Specifies the cloud resource IP.

* `ip_list` - (Optional, List) Specifies the cloud resource IP list.

* `resource_id_list` - (Optional, List) Specifies the resource ID list.

* `status` - (Optional, String) Specifies the resource status.

  For details, see [status](https://support.huaweicloud.com/api-ecs/ecs_08_0002.html)

* `agent_state` - (Optional, String) Specifies the unified agent status.

* `image_name` - (Optional, String) Specifies the image name.

* `os_type` - (Optional, String) Specifies the cloud resource operating system type.

* `tag` - (Optional, String) Specifies the tags for cloud resources. The format of the tag is **key.value**.
  When naming tags, the following requirements must be met:
  + The **key** of the tag can only contain uppercase letters (A~Z), lowercase letters (a~z), numbers (0-9),
    underscores (\_), hyphens (-), and Chinese characters.
  + The **value** of the tag can only contain uppercase letters (A~Z),lowercase letters (a~z), numbers (0-9),
    underscores (\_), hyphens (-), decimal points (.), and Chinese characters.

* `tag_key` - (Optional, String) Specifies the tag key of the cloud resource.

* `group_id` - (Optional, String) Specifies the group ID of the cloud resource.

* `component_id` - (Optional, String) Specifies the component ID of the cloud resource.

* `application_id` - (Optional, String) Specifies the application ID of the cloud resource.

* `cce_cluster_id` - (Optional, String) Specifies the CCE cluster ID.

* `vpc_id` - (Optional, String) Specifies the virtual private cloud ID.

* `is_delegated` - (Optional, Bool) Specifies whether the resource is delegated.

* `operable` - (Optional, String) Specifies whether the user defined resource can operate the instance. If the value
  is **enable**, it is enabled; if the current field does not exist, it is not enabled.

* `is_collected` - (Optional, Bool) Specifies whether it is a favorite.

* `flavor_name` - (Optional, String) Specifies the cloud resource specification name.

* `charging_mode` - (Optional, String) Specifies the billing type for the cloud server.
  Values can be as follows:
  + **0**: On demand billing.
  + **1**: Yearly package or monthly package.
  + **2**: Award based billing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the resource list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the resource ID assigned by the CMDB.

* `resource_id` - Indicates the resource ID.

* `name` - Indicates the resource name.

* `ep_id` - Indicates the enterprise project ID.

* `project_id` - Indicates the project ID in OpenStack.

* `domain_id` - Indicates the tenant ID.

* `cloud_service_name` - Indicates the cloud service name.

* `type` - Indicates the resource type.

* `region_id` - Indicates the region ID.

* `tags` - Indicates the resource tags.

  The [tags](#data_tags_struct) structure is documented below.

* `properties` - Indicates the resource properties.

* `ingest_properties` - Indicates the ingest attributes of the resource.

* `agent_id` - Indicates the ID assigned by unified agent.

* `agent_state` - Indicates the unified agent status.

* `is_delegated` - Indicates whether the resource is delegated.

* `operable` - Indicates whether the user defined resource can operate the instance. If the value is **enable**, it is
  enabled; if the current field does not exist, it is not enabled.

<a name="data_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.
