# Package antiddos
    import "github.com/huaweicloud/golangsdk/openstack/antiddos/v1/antiddos"
**[概述](#概述)**  

**[目录](#目录)**  

**[API对应表](#API对应表)**  

**[开始](#开始)**  

## W概述
Anti-DDoS流量清洗（Anti-DDoS Service）是通过专业的防DDoS设备来为客户互联网应用提供精细化的抵御DDoS攻击能力（包括CC、SYN flood、UDP flood等所有DDoS攻击方式）。可根据租用带宽及业务模型自助配置防护阈值参数，系统检测到攻击后通知用户进行网站防御。

示例代码，更新指定EIP的Anti-DDoS防护策略配置。

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
示例代码， 用户开通Anti-DDoS流量清洗防护。

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Create(client.ServiceClient(), floatingIpId, createOpt).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 查询指定EIP在过去24小时之内的防护流量信息，流量的间隔时间单位为5分钟。

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.DailyReport(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 用户关闭Anti-DDoS流量清洗防护。

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Delete(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 查询配置的Anti-DDoS防护策略。

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.Get(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 查询指定EIP的Anti-DDoS防护状态。

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.GetStatus(client.ServiceClient(), floatingIpId).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 用户查询指定的Anti-DDoS防护配置任务。

    actual, err := antiddos.GetTask(client.ServiceClient(), antiddos.GetTaskOpts{
        TaskId: "4a4fefe7-34a1-40e2-a87c-16932af3ac4a",
    }).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 查询系统支持的Anti-DDoS防护策略配置的可选范围。

    actual, err := antiddos.ListConfigs(client.ServiceClient()).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 查询指定EIP在过去24小时之内的异常事件信息。

    floatingIpId := "82abaa86-8518-47db-8d63-ddf152824635"
    actual, err := antiddos.ListLogs(client.ServiceClient(), floatingIpId, antiddos.ListLogsOpts{
        Limit:   2,
        Offset:  1,
        SortDir: "asc",
    }).Extract()
    if err != nil {
      panic(err)
    }
示例代码， 查询用户所有EIP的Anti-DDoS防护状态信息。

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
示例代码， 查询用户所有Anti-DDoS防护周统计情况。

    actual, err := antiddos.WeeklyReport(client.ServiceClient(), antiddos.WeeklyReportOpts{}).Extract()
    if err != nil {
      panic(err)
    }
## 目录
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
## API对应表
|类别|API|EndPoint|
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
## 开始
## func Create
    func Create(*golangsdk.ServiceClient, string, CreateOptsBuilder) (CreateResult)  
用户开通Anti-DDoS流量清洗防护。作为异步接口，调用成功，只是说明服务节点收到了开通请求，开通是否成功需要通过任务查询接口查询该任务的执行状。
## func DailyReport
    func DailyReport(*golangsdk.ServiceClient, string) (DailyReportResult)  
查询指定EIP在过去24小时之内的防护流量信息，流量的间隔时间单位为5分钟。
## func Delete
    func Delete(*golangsdk.ServiceClient, string) (DeleteResult)  
用户关闭Anti-DDoS流量清洗防护。作为异步接口，调用成功，只是说明服务节点收到了关闭防护请求，操作是否成功需要通过任务查询接口查询该任务的执行状态。
## func Get
    func Get(*golangsdk.ServiceClient, string) (GetResult)  
查询配置的Anti-DDoS防护策略，用户可以查询指定EIP的Anti-DDoS防护策略。
## func GetStatus
    func GetStatus(*golangsdk.ServiceClient, string) (GetStatusResult)  
查询指定EIP的Anti-DDoS防护状态。
## func GetTask
    func GetTask(*golangsdk.ServiceClient, GetTaskOptsBuilder) (GetTaskResult)  
用户查询指定的Anti-DDoS防护配置任务，得到任务当前执行的状态。
## func ListConfigs
    func ListConfigs(*golangsdk.ServiceClient) (ListConfigsResult)  
查询系统支持的Anti-DDoS防护策略配置的可选范围，用户根据范围列表选择适合自已业务的防护策略进行Anti-DDoS流量清洗。
## func ListLogs
    func ListLogs(*golangsdk.ServiceClient, string, ListLogsOptsBuilder) (ListLogsResult)  
查询指定EIP在过去24小时之内的异常事件信息，异常事件包括清洗事件和黑洞事件，查询延迟在5分钟之内。
## func ListStatus
    func ListStatus(*golangsdk.ServiceClient, ListStatusOptsBuilder) (ListStatusResult)  
查询用户所有EIP的Anti-DDoS防护状态信息，用户的EIP无论是否绑定到云服务器，都可以进行查询。
## func Update
    func Update(*golangsdk.ServiceClient, string, UpdateOptsBuilder) (UpdateResult)  
更新指定EIP的Anti-DDoS防护策略配置。调用成功，只是说明服务节点收到了关闭更新配置请求，操作是否成功需要通过任务查询接口查询该任务的执行状态。
## func WeeklyReport
    func WeeklyReport(*golangsdk.ServiceClient, WeeklyReportOptsBuilder) (WeeklyReportResult)  
查询用户所有Anti-DDoS防护周统计情况，包括一周内DDoS拦截次数和攻击次数、以及按照被攻击次数进行的排名信息等统计数据。系统支持当前时间之前四周的周统计数据查询，超过这个时间的请求是查询不到统计数据的。
