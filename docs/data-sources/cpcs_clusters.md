---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_clusters"
description: |-
  Use this data source to get the list of CPCS clusters.
---

# huaweicloud_cpcs_clusters

Use this data source to get the list of CPCS clusters.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
data "huaweicloud_cpcs_clusters" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the cluster.

* `service_type` - (Optional, String) Specifies the service type of the cluster.
  Valid values are:
  + **ENCRYPT_DECRYPT**: Encryption and decryption services.
  + **SIGN_VERIFY**: Signature verification services.
  + **KMS**: Key Management Service (KMS).
  + **TIMESTAMP**: Timestamp services.
  + **COLLA_SIGN**: Collaborative signature services.
  + **OTP**: Dynamic password services.
  + **DB_ENCRYPT**: Database encryption services.
  + **FILE_ENCRYPT**: File encryption services.
  + **DIGIT_SEAL**: Electronic signature and seal services.
  + **SSL_VPN**: SSL and VPN services.

* `sort_key` - (Optional, String) Specifies the sort attribute.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The default value is **DESC**. Valid values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - Indicates the clusters list.
  The [clusters](#CPCS_clusters) structure is documented below.

<a name="CPCS_clusters"></a>
The `clusters` block supports:

* `task_id` - The task ID.

* `project_id` - The project ID.

* `domain_id` - The account ID.

* `ccsp_id` - The CCSP cluster ID.

* `distributed_type` - The distribution type of the cluster.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `service_type` - The service type of the cluster.

* `type` - The type of the cluster. Valid values are `SHARED` and `EXCLUSIVE`.

* `instance_num` - The number of service instances in the cluster.

* `status` - The status of the cluster.

* `progress_info` - The progress information.

* `vsm_num` - The number of VSM instances used by the cluster.

* `create_time` - The creation time of the cluster, in UNIX timestamp format.

* `shared_ccsp` - Whether the CCSP is shared.

* `enterprise_project_id` - The enterprise project ID.

* `az` - The availability zone where the cluster is located.
