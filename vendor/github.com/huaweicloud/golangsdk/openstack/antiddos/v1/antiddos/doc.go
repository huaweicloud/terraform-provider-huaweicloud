/*
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
*/
package antiddos
