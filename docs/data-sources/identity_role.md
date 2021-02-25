---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_role

Use this data source to get the ID of an HuaweiCloud role.

The Role in Terraform is the same as Policy on console. however,
The policy name is the display name of Role, the Role name cannot
be found on Console. please refer to the following table to configuration
Role:

Role Name | Policy Name
---- | ---
system_all_0 | All permissions of ECS service
system_all_1 | Common permissions of ECS service, except installation, delete, reinstallation and so on
system_all_2 | The read-only permissions to all ECS resources, which can be used for statistics and survey
system_all_3 | The read-only permissions to all EVS resources, which can be used for statistics and survey
system_all_4 | All permissions of EVS service
system_all_5 | All permissions of VPC service
system_all_6 | The read-only permissions to all VPC resources, which can be used for statistics and survey
system_all_7 | The read-only permissions to all Document Database Service resources, which can be used for statistics and survey
system_all_8 | DBA permissions of Document Database Service, except delete
system_all_9 | All permissions of Document Database Service
system_all_10 | DBA permissions of Relational Database Service, except delete
system_all_11 | All permissions of Relational Database Service
system_all_12 | The read-only permissions to all Relational Database Service resources, which can be used for statistics and survey
system_all_13 | Cloud Container Engine Cluster Viewer
system_all_14 | Cloud Container Engine Cluster Admin
system_all_1001 | Full access to all resources
secu_admin | Security Administrator
te_admin | Tenant Administrator
te_agency | Agent Operator
readonly | Guest
server_adm | Server Administrator
as_adm | AutoScaling Administrator
aos_adm | Application Orchestration Service Admin
aos_dev | Application Orchestration Service Developer
vbs_adm | Volume Backup Service Administrator
tms_adm | Tag Management Service Administrator
dcs_admin | Distributed Cache Service Administrator
swr_adm | Software Repository Admin
elb_adm | ELB Service Administrator
dss_adm | Dedicated Storage Service Administrator
dws_adm | Data Warehouse Service Administrator
kms_adm | KMS Administrator
ims_adm | IMS Administrator
ddos_adm | Anti-DDoS Administrator
dns_adm | DNS Administrator
wks_adm | Workspace Administrator
nat_adm | NAT Gateway Administrator
cse_adm | Cloud Service Engine Admin
rds_adm | RDS Administrator
dis_adm | DIS Administrator
sfs_adm | Scalable File Service Administrator
smn_adm | SMN Administrator
cts_adm | CTS Administrator
apm_adm | Application Performance Monitor Admin
mrs_adm | MRS Administrator
ces_adm | CES Administrator
rts_adm | RTS Service Administrator
cce_adm | CCE Administrator
cs_adm | Cloud Stream Service Tenant Administrator, can manage multiple CS users
cs_user | Cloud Stream Service User
dms_adm | DMS Administrator
dps_adm | DPS Administrator
mls_adm | Machine Learning Service Administrator
css_adm | Cloud Search Service Administrator
dds_adm | Document Database Service Administrator
csbs_adm | Cloud Server Backup Service Administrator
sdrs_adm | Storage Disaster Recovery Service Administrator
svcstg_adm | ServiceStage Admin
svcstg_dev | ServiceStage Developer
svcstg_opr | ServiceStage Operator
vpc_netadm | VPC Administrator
vpcep_adm | VPCEndpoint service enables you to privately connect your VPC to supported services

```hcl
data "huaweicloud_identity_role" "auth_admin" {
  name = "secu_admin"
}
```

## Argument Reference

* `name` - (Required, String) The name of the role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.
