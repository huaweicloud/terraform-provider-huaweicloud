 package main

// import (
//     "fmt"
//     "huaweicloud-sdk-go-v3/core/auth/basic"
//     sms "huaweicloud-sdk-go-v3/services/sms/v3"
// 	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/model"
//     region "huaweicloud-sdk-go-v3/services/sms/v3/region"
// )

// func main() {
//     ak := "OYFWTNGHT3HG0V0MTXBZ"
//     sk := "TYkhq0NqrCmOjAckdVC4fhnCW0AgGdJFAZ2q1oCQ"

//     // auth := baisc.NewCredentialsBuilder().
//     //     WithAk(ak).
//     //     WithSk(sk).
//     //     Build()
//     auth:=basic.NewCredentialsBuilder().WithAk(ak).WithSk(sk).Build()

//     client := sms.NewSmsClient(
//         sms.SmsClientBuilder().
//             WithRegion(region.ValueOf("ap-southeast-1")).
//             WithCredential(auth).
//             Build())

//     request := &model.CreateTemplateRequest{}
// 	availabilityZoneTemplateTemplateRequest:= "ap-southeast-1a"
// 	templatebody := &model.TemplateRequest{
// 		Name: "test8",
// 		IsTemplate: true,
// 		Region: "ap-southeast-1",
// 		Projectid: "0e1b256e9680f3722fdec005b172cc1f",
// 		AvailabilityZone: &availabilityZoneTemplateTemplateRequest,
// 	}
// 	request.Body = &model.CreateTemplateReq{
// 		Template: templatebody,
// 	}
// 	response, err := client.CreateTemplate(request)
// 	if err == nil {
//         fmt.Printf("%+v\n", response)
//     } else {
//         fmt.Println(err)
//     }
// }