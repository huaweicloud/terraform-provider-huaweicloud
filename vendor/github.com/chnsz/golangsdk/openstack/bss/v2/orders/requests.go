package orders

import (
	"github.com/chnsz/golangsdk"
)

type UnsubscribeOpts struct {
	ResourceIds     []string `json:"resource_ids" required:"true"`
	UnsubscribeType int      `json:"unsubscribe_type" required:"true"`
}

type UnsubscribeOptsBuilder interface {
	ToOrderUnsubscribeMap() (map[string]interface{}, error)
}

func (opts UnsubscribeOpts) ToOrderUnsubscribeMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Unsubscribe(client *golangsdk.ServiceClient, opts UnsubscribeOptsBuilder) (r UnsubscribeResult) {
	reqBody, err := opts.ToOrderUnsubscribeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(unsubscribeURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

type PayOrderOpts struct {
	OrderId string `json:"order_id" required:"true"`
	// Whether coupons are used for order payment. If this parameter is set to YES, the coupon_infos
	// field is mandatory. If this parameter is set to NO, the value of coupon_infos is ignored.
	UseCoupon string `json:"use_coupon" required:"true"`
	// Whether a discount is used for order payment. If this parameter is set to YES, the discount_infos
	// field is mandatory. If this parameter is set to NO, the value of discount_infos is ignored.
	UseDiscount string `json:"use_discount" required:"true"`
	// The coupons info, if UseCoupon is set to YES, this parameter is mandatory
	CouponInfos []CouponSimpleInfo `json:"coupon_infos,omitempty"`
	// The discounts info, if UseDiscount is set to YES, this parameter is mandatory
	DiscountInfos []DiscountSimpleInfo `json:"discount_infos,omitempty"`
}

type CouponSimpleInfo struct {
	// Coupon ID
	Id string `json:"id" required:"true"`
	// Coupon type
	// 300: Discount coupon (reserved)
	// 301: Promotion coupon
	// 302: Promotion flexi-purchase coupon (reserved)
	// 303: Promotion stored-value card (reserved)
	Type int `json:"type" required:"true"`
}

type DiscountSimpleInfo struct {
	// Discount ID
	Id string `json:"id" required:"true"`
	// Discount type
	// 0: promotion discount
	// 2: commercial discount
	// 3: discount granted by a partner
	Type int `json:"type" required:"true"`
}

type PayOrderOptsBuilder interface {
	ToPayOrderMap() (map[string]interface{}, error)
}

func (opts PayOrderOpts) ToPayOrderMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func PayOrder(client *golangsdk.ServiceClient, opts PayOrderOptsBuilder) (r PayOrderResult) {
	reqBody, err := opts.ToPayOrderMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(payOrderURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{204}})
	return
}
