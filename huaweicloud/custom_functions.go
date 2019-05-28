package huaweicloud

func expandCdmClusterV1CreateAutoRemind(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	email, _ := navigateValue(d, []string{"email"}, nil)
	e, eok := email.([]interface{})

	phone, _ := navigateValue(d, []string{"phone_num"}, nil)
	p, pok := phone.([]interface{})

	return (eok && len(e) > 0) || (pok && len(p) > 0), nil
}

func expandCdmClusterV1CreateClusterIsScheduleBootOff(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	on, _ := navigateValue(d, []string{"schedule_boot_time"}, nil)
	on1, ok1 := on.(string)

	off, _ := navigateValue(d, []string{"schedule_off_time"}, nil)
	off1, ok2 := off.(string)

	return (ok1 && on1 != "") || (ok2 && off1 != ""), nil
}
