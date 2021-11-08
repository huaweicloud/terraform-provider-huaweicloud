package main

// import (
// 	"fmt"
// 	"huaweicloud-sdk-go-v3/core/auth/basic"
//     sms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3"
// 	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/model"
//     region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/region"
// )

// func main() {
// 	ak := "OYFWTNGHT3HG0V0MTXBZ"
// 	sk := "TYkhq0NqrCmOjAckdVC4fhnCW0AgGdJFAZ2q1oCQ"

//     auth := basic.NewCredentialsBuilder().
//         WithAk(ak).
//         WithSk(sk).
//         Build()

//     client := sms.NewSmsClient(
//         sms.SmsClientBuilder().
//             WithRegion(region.ValueOf("ap-southeast-1")).
//             WithCredential(auth).
//             Build())

//     request := &model.DeleteTemplatesRequest{}
// 	var listIdsbody = []string{
//         "35056fca-d7f9-4f67-b9c9-4d11c80d83c5",
//     }
// 	request.Body = &model.DeletetemplatesReq{
// 		Ids: &listIdsbody,
// 	}
// 	response, err := client.DeleteTemplates(request)
// 	if err == nil {
//         fmt.Printf("%+v\n", response)
//     } else {
//         fmt.Println(err)
//     }
// }