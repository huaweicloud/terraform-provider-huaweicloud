# Package antiddos
    import "github.com/huaweicloud/golangsdk/openstack/antiddos/v1/antiddos"
**[Overview](#overview)**  

**[Index](#index)**  

**[API Mapping](#api-mapping)**  

**[Content](#content)**  

## Overview
The Anti-DDoS traffic cleaning service (Anti-DDoS for short) defends resources (Elastic Cloud Servers (ECSs), Elastic Load Balance (ELB) instances, and Bare Metal Servers (BMSs)) on HUAWEI CLOUD against network- and application-layer distributed denial of service (DDoS) attacks and sends alarms immediately when detecting an attack. In addition, Anti-DDoS improves the utilization of bandwidth and ensures the stable running of users' services.

Example to update the Anti-DDoS defense policy of a specified EIP.

    updateOpt := antiddos.UpdateOpts{
      EnableL7:            true,
      TrafficPosId:        1,
      HttpRequestPosId:    2,
      CleaningAccessPosId: 3,
      AppTypeId:           1,
    }
    
    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Update(client.ServiceClient(), floatingIpId, updateOpt).Extract()
    if err != nil {
      panic(err)
    }
Example to enable the Anti-DDoS traffic cleaning defense. 

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Create(client.ServiceClient(), floatingIpId, createOpt).Extract()
    if err != nil {
      panic(err)
    }
Example to query the traffic of a specified EIP in the last 24 hours. Traffic is detected in five-minute intervals.

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.DailyReport(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
Example to disable the Anti-DDoS traffic cleaning defense.

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Delete(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
Example to query configured Anti-DDoS defense policies.

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Get(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
Example to query the defense status of a specified EIP.

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.GetStatus(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
Example to query the execution status of a specified Anti-DDoS configuration task.

    actual, err := antiddos.GetTask(client.ServiceClient(), antiddos.GetTaskOpts{
        TaskId: "4a4fefe7-34a1-40e2-a87c-16932af3ac4a",
    }).Extract()
    if err != nil {
      panic(err)
    }
Example to query optional Anti-DDoS defense policies. 

    actual, err := antiddos.ListConfigs(client.ServiceClient()).Extract()
    if err != nil {
      panic(err)
    }
Example to query events of a specified EIP in the last 24 hours.

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.ListLogs(client.ServiceClient(), floatingIpId, antiddos.ListLogsOpts{
        Limit:   2,
        Offset:  1,
        SortDir: "asc",
    }).Extract()
    if err != nil {
      panic(err)
    }
Example to query the defense statuses of all EIPs.

    listOpt := antiddos.ListStatusOpts{
        Limit:  2,
        Offset: 1,
        Status: "notConfig",
        Ip:     "49.",
    }
    
    actual, err := antiddos.ListStatus(client.ServiceClient(), listOpt).Extract()
    if err != nil {
      panic(err)
    }
Example to query weekly defense statistics about all your EIPs.

    actual, err := antiddos.WeeklyReport(client.ServiceClient(), antiddos.WeeklyReportOpts{}).Extract()
    if err != nil {
      panic(err)
    }
## Index
**[func Create(*golangsdk.ServiceClient, string, CreateOptsBuilder) (CreateResult)](#func-create)**  
**[func DailyReport(*golangsdk.ServiceClient, string) (DailyReportResult)](#func-dailyreport)**  
**[func Delete(*golangsdk.ServiceClient, string) (DeleteResult)](#func-delete)**  
**[func Get(*golangsdk.ServiceClient, string) (GetResult)](#func-get)**  
**[func GetStatus(*golangsdk.ServiceClient, string) (GetStatusResult)](#func-getstatus)**  
**[func GetTask(*golangsdk.ServiceClient, GetTaskOptsBuilder) (GetTaskResult)](#func-gettask)**  
**[func ListConfigs(*golangsdk.ServiceClient) (ListConfigsResult)](#func-listconfigs)**  
**[func ListLogs(*golangsdk.ServiceClient, string, ListLogsOptsBuilder) (ListLogsResult)](#func-listlogs)**  
**[func ListStatus(*golangsdk.ServiceClient, ListStatusOptsBuilder) (ListStatusResult)](#func-liststatus)**  
**[func Update(*golangsdk.ServiceClient, string, UpdateOptsBuilder) (UpdateResult)](#func-update)**  
**[func WeeklyReport(*golangsdk.ServiceClient, WeeklyReportOptsBuilder) (WeeklyReportResult)](#func-weeklyreport)**  
## API Mapping
|Catalog|API|EndPoint|
|----|---|--------|
|antiddos|func Create(*golangsdk.ServiceClient, string, CreateOptsBuilder) (CreateResult)|POST /v1/{project_id}/antiddos/{floating_ip_id}|
|antiddos|func DailyReport(*golangsdk.ServiceClient, string) (DailyReportResult)|GET /v1/{project_id}/antiddos/{floating_ip_id}/daily|
|antiddos|func Delete(*golangsdk.ServiceClient, string) (DeleteResult)|DELETE /v1/{project_id}/antiddos/{floating_ip_id}|
|antiddos|func Get(*golangsdk.ServiceClient, string) (GetResult)|GET /v1/{project_id}/antiddos/{floating_ip_id}|
|antiddos|func GetStatus(*golangsdk.ServiceClient, string) (GetStatusResult)|GET /v1/{project_id}/antiddos/{floating_ip_id}/status|
|antiddos|func GetTask(*golangsdk.ServiceClient, GetTaskOptsBuilder) (GetTaskResult)|GET /v1/{project_id}/query_task_status|
|antiddos|func ListConfigs(*golangsdk.ServiceClient) (ListConfigsResult)|GET /v1/{project_id}/antiddos/query_config_list|
|antiddos|func ListLogs(*golangsdk.ServiceClient, string, ListLogsOptsBuilder) (ListLogsResult)|GET /v1/{project_id}/antiddos/{floating_ip_id}/logs|
|antiddos|func ListStatus(*golangsdk.ServiceClient, ListStatusOptsBuilder) (ListStatusResult)|GET /v1/{project_id}/antiddos|
|antiddos|func Update(*golangsdk.ServiceClient, string, UpdateOptsBuilder) (UpdateResult)|PUT /v1/{project_id}/antiddos/{floating_ip_id}|
|antiddos|func WeeklyReport(*golangsdk.ServiceClient, WeeklyReportOptsBuilder) (WeeklyReportResult)|GET /v1/{project_id}/antiddos/weekly|
## Content
## func Create
    func Create(*golangsdk.ServiceClient, string, CreateOptsBuilder) (CreateResult)  
This asynchronous API allows you to enable the Anti-DDoS traffic cleaning defense. Successfully invoking this API only means that the service node has received the enabling request. You need to use the task querying API to check the task execution status. 
## func DailyReport
    func DailyReport(*golangsdk.ServiceClient, string) (DailyReportResult)  
This API allows you to query the traffic of a specified EIP in the last 24 hours. Traffic is detected in five-minute intervals.
## func Delete
    func Delete(*golangsdk.ServiceClient, string) (DeleteResult)  
This asynchronous API allows you to disable the Anti-DDoS traffic cleaning defense. Successfully invoking this API only means that the service node has received the disabling request. You need to use the task querying API to check the task execution status.
## func Get
    func Get(*golangsdk.ServiceClient, string) (GetResult)  
This API enables you to query configured Anti-DDoS defense policies. You can query the policy of a specified EIP.
## func GetStatus
    func GetStatus(*golangsdk.ServiceClient, string) (GetStatusResult)  
This API allows you to query the defense status of a specified EIP.
## func GetTask
    func GetTask(*golangsdk.ServiceClient, GetTaskOptsBuilder) (GetTaskResult)  
This API enables you to query the execution status of a specified Anti-DDoS configuration task.
## func ListConfigs
    func ListConfigs(*golangsdk.ServiceClient) (ListConfigsResult)  
This API allows you to query optional Anti-DDoS defense policies. Based on your service, you can select a policy for Anti-DDoS traffic cleaning.
## func ListLogs
    func ListLogs(*golangsdk.ServiceClient, string, ListLogsOptsBuilder) (ListLogsResult)  
This API allows you to query events of a specified EIP in the last 24 hours. Events include cleaning and blackhole events, and the query delay is within five minutes.
## func ListStatus
    func ListStatus(*golangsdk.ServiceClient, ListStatusOptsBuilder) (ListStatusResult)  
This API enables you to query the defense statuses of all EIPs, regardless whether an EIP has been bound to an Elastic Cloud Server (ECS) or not.
## func Update
    func Update(*golangsdk.ServiceClient, string, UpdateOptsBuilder) (UpdateResult)  
This API enables you to update the Anti-DDoS defense policy of a specified EIP. Successfully invoking this API only means that the service node has received the update request. You need to use the task querying API to check the task execution status. 
## func WeeklyReport
    func WeeklyReport(*golangsdk.ServiceClient, WeeklyReportOptsBuilder) (WeeklyReportResult)  
This API allows you to query weekly defense statistics about all your EIPs, including the number of intercepted DDoS attacks, number of attacks, and ranking by the number of attacks. Currently, you can query weekly statistics up to four weeks before the current time. Data older than four weeks cannot be queried.
