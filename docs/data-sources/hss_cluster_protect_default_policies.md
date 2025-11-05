---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_protect_default_policies"
description: |-
  Use this data source to get the list of HSS cluster protect default policies within HuaweiCloud.
---

# huaweicloud_hss_cluster_protect_default_policies

Use this data source to get the list of HSS cluster protect default policies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cluster_protect_default_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `general_policy_num` - The number of general policies.

* `malicious_image_policy_num` - The number of malicious image policies.

* `security_policy_num` - The number of security policies.

* `data_list` - The list of cluster protect default policies.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `content` - The policy content.

* `deploy_content` - The deployment content.

* `parameters` - The parameters.

* `policy_name` - The policy name.

* `policy_id` - The policy ID.

* `resources` - The resources.

  The [resources](#data_list_resources_struct) structure is documented below.

* `template_id` - The template ID.

* `template_name` - The template name.

* `template_type` - The template type.

* `description` - The policy description.

* `update_time` - The update time.

* `create_time` - The creation time.

* `image_num` - The number of protected images.

* `labels_num` - The number of protected labels.

* `status` - The status.

* `white_images` - The white list images.

  The [white_images](#data_list_white_images_struct) structure is documented below.

<a name="data_list_resources_struct"></a>
The `resources` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `images` - The images.

* `labels` - The labels.

* `namespace` - The namespace.

<a name="data_list_white_images_struct"></a>
The `white_images` block supports:

* `cluster_id` - The cluster ID.

* `image_name` - The image name.

* `image_version` - The image version.
