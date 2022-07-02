package conferences

import (
	"github.com/chnsz/golangsdk"
)

type CreateOpts struct {
	// The conference start time (UTC time).
	//   When creating a reservation conference, if the start time is not specified or the blank string is filled in, it
	//   means that the conference will start immediately.
	//   The time is in UTC, which is the time in time zone 0.
	// Format: yyyy-MM-dd HH:mm
	StartTime string `json:"startTime,omitempty"`
	// The duration of the conference, in minutes, with a maximum value of 1440 and a minimum value of 15.
	// The default is 30.
	Length int `json:"length,omitempty"`
	// Conference subject. The length is limited to 128 characters.
	Subject string `json:"subject,omitempty"`
	// The media type of the conference.
	// It consists of one or more enumeration Strings. When there are multiple enumerations, each enumeration value is
	// separated by "," commas. The enumeration values are as follows:
	//   "Voice": Voice.
	//   "Video": SD video.
	//   "HDVideo": High-definition video (mutually exclusive with Video, if Video and HDVideo are selected at the same
	//     time, the system will select Video by default).
	//   "Telepresence": Telepresence (mutually exclusive with HDVideo and Video, if selected at the same time, the
	//     system uses Telepresence). (reserved field)
	//   "Data": Multimedia (system configuration determines whether to automatically add Data).
	MediaTypes string `json:"mediaTypes,omitempty"`
	// When the soft terminal creates an instant conference, it carries a temporary group ID in the current field, which
	// is carried by the server in the conference-info header field when inviting other participants.
	// The length is limited to 31 characters.
	Groupuri string `json:"groupuri,omitempty"`
	// Attendee list. This list can be used to send conference notifications, conference reminders, and automatic invitations
	// when the conference starts.
	Participants []Participant `json:"attendees,omitempty"`
	// The parameter of the cyclical conference. When the conference is a periodic conference, this parameter must be filled in,
	// otherwise the server ignores this parameter.
	// This parameter includes the start date, end date of the periodic conference, the period of the conference and the
	// conference time point in the period.
	CycleParams *CycleParams `json:"cycleParams,omitempty"`
	// Whether the conference automatically starts recording, it only takes effect when the recording type is:
	//   1: Automatically start recording.
	//   0: Do not start recording automatically.
	// The default is not to start automatically.
	IsAutoRecord *int `json:"isAutoRecord,omitempty"`
	// Conference media encryption mode.
	//   0: Adaptive encryption.
	//   1: Force encryption.
	//   2: Do not encrypt.
	// Default values are populated by enterprise-level configuration.
	EncryptMode *int `json:"encryptMode,omitempty"`
	// The default language of the conference, the default value is defined by the conference cloud service.
	// For languages supported by the system, it is passed according to the RFC3066 specification.
	//   zh-CN: Simplified Chinese.
	//   en-US: US English.
	Language string `json:"language,omitempty"`
	// The time zone information of the conference time in the conference notification.
	// For time zone information, refer to the time zone mapping relationship.
	// For example: "timeZoneID":"26", the time in the conference notification sent through HUAWEI CLOUD conference will be
	// marked as "2021/11/11 Thursday 00:00 - 02:00 (GMT) Greenwich Standard When: Dublin, Edinburgh, Lisbon, London".
	// For an aperiodic conference, if the conference notification is sent through a third-party system, this field does not
	// need to be filled in.
	TimeZoneID string `json:"timeZoneID,omitempty"`
	// Recording type.
	//   0: Disabled.
	//   1: Live broadcast.
	//   2: Record and broadcast.
	//   3: Live + Recording.
	// Default is disabled.
	RecordType *int `json:"recordType,omitempty"`
	// The mainstream live broadcast address, with a maximum of 255 characters.
	// It is valid when the recording type is 2 or 3.
	LiveAddress string `json:"liveAddress,omitempty"`
	// Auxiliary streaming address, the maximum length is 255 characters.
	// It is valid when the recording type is 2 or 3.
	AuxAddress string `json:"auxAddress,omitempty"`
	// Whether to record auxiliary stream.
	//   0: Do not record.
	//   1: Record.
	// It is valid when the recording type is 2 or 3.
	RecordAuxStream *int `json:"recordAuxStream,omitempty"`
	// Other configuration information for the conference.
	Configuration *Configuration `json:"confConfigInfo,omitempty"`
	// Recording authentication method.
	//   0: Viewable/downloadable via link.
	//   1: Enterprise users can watch/download.
	//   2: Attendees can watch/download.
	// It is valid when the recording type is 2 or 3.
	RecordAuthType *int `json:"recordAuthType,omitempty"`
	// Whether to use the cloud conference room or personal conference ID to hold a reservation conference.
	//   0: Do not use cloud conference room.
	//   1: Use cloud conference room or personal conference ID.
	// Cloud conference rooms are not used by default.
	VmrFlag *int `json:"vmrFlag,omitempty"`
	// VMR ID bound to the current founding account.
	// Obtained by querying the cloud conference room and personal conference ID interface.
	// Note:
	//   vmrID takes the ID returned in the above query interface, not vmrId.
	//   When creating a conference with a personal conference ID, use VMR with vmrMode=0;
	//   When creating a conference in a cloud conference room, use VMR with vmrMode=1.
	VmrId string `json:"vmrID,omitempty"`
	// The number of parties in the conference, the maximum number of participants in the conference.
	//   0: Unlimited.
	// Greater than 0: the maximum number of participants in the conference.
	ConcurrentParticipants *int `json:"concurrentParticipants,omitempty"`
	// The authorization token.
	Token string `json:"-" required:"true"`
}

type Participant struct {
	// User UUID of the participant.
	UserUUID string `json:"userUUID,omitempty"`
	// The account ID of the participant.
	// If it is an account/password authentication scenario, optional, indicating the ID of the HUAWEI CLOUD conference
	// account. If it is an APP ID authentication scenario, it is required, indicating the user ID of the third party,
	// and the appid parameter needs to be carried.
	AccountId string `json:"accountId,omitempty"`
	// App ID, application identification, an application only needs to be created once, refer to "App ID application"
	// If it is an APP ID authentication scenario, this item is required.
	AppId string `json:"appId,omitempty"`
	// Attendee name or nickname. The length is limited to 96 characters.
	Name string `json:"name,omitempty"`
	// The role in the conference. The default is a regular participant.
	//   0: Normal attendees.
	//   1: The conference host.
	Role *int `json:"role,omitempty"`
	// If it is an account/password authentication scenario, it is required to fill in the number (SIP and TEL number
	// formats can be supported). Maximum of 127 characters. At least one of phone, email and sms must be filled in.
	// If it is an APP ID authentication scenario, optional.
	Phone string `json:"phone,omitempty"`
	// Email address. Maximum of 255 characters. In the case of account/password authentication, at least one of phone,
	// email, and sms must be filled in. (Information notification for reservation, modification and cancellation of
	// conference)
	Email string `json:"email,omitempty"`
	// Mobile number for SMS notification. Maximum of 32 characters. In the case of account/password authentication, at
	// least one of phone, email, and sms must be filled in. (Information notification for reservation, modification and
	// cancellation of conference)
	SMS string `json:"sms,omitempty"`
	// Whether the user needs to be automatically muted when joining the conference (only effective when invited in the
	// conference). Unmuted by default.
	//   0: No mute required.
	//   1: Mute is required.
	IsMute *int `json:"isMute,omitempty"`
	// Whether to automatically invite this participant when the conference starts. The default value is determined by the
	// enterprise-level configuration.
	//   0: Do not automatically invite.
	//   1: Automatic invitation.
	IsAutoInvite *int `json:"isAutoInvite,omitempty"`
	// The default value is defined by the conference AS. The number type enumeration is as follows:
	//   normal: Soft terminal.
	//   telepresence: Telepresence. Single-screen and triple-screen telepresence belong to this category. (reserved field)
	//   terminal: conference room or hard terminal.
	//   outside: The outside participant.
	//   mobile: User's mobile phone number.
	//   telephone: The user's landline phone. (reserved field)
	//   ideahub: ideahub.
	Type string `json:"type,omitempty"`
	// Department ID. Maximum of 64 characters.
	DeptUUID string `json:"deptUUID,omitempty"`
	// Department name. Maximum of 128 characters.
	DeptName string `json:"deptName,omitempty"`
}

type CycleParams struct {
	// The start date of the recurring conference, format: YYYY-MM-DD.
	// The start date cannot be earlier than the current date.
	// The date is the date in the time zone specified by timeZoneID, not the date in UTC time.
	StartDate string `json:"startDate,omitempty"`
	// End date of the recurring conference, format: YYYY-MM-DD.
	// The maximum time interval between start date and end date cannot exceed 1 year.
	// A maximum of 50 subconferences are allowed between the start date and the end date.
	// If there are more than 50 subconferences, the end date will be automatically adjusted.
	// The date is the date in the time zone specified by timeZoneID, not the date in UTC time.
	EndDate string `json:"endDate,omitempty"`
	// Period type. The valid values are:
	//   Day
	//   Week
	//   Month
	Cycle string `json:"cycle,omitempty"`
	// If the cycle selects "Day", which means that it will be held every few days, and the value range is [1,15].
	// If the cycle selects "Week", which means that it is held every few weeks, and the value range is [1,5].
	// If the cycle selects "Month", Interval means every few months, the value range is [1,3].
	Interval int `json:"interval,omitempty"`
	// The conference point in the cycle. Only valid by week and month.
	// "Week" is selected for "cycle", and two elements 1 and 3 are filled in point, which means that a conference is held
	// every Monday and Wednesday, and 0 means Sunday.
	// "Month" is selected for "cycle", and 12 and 20 are filled in point, which means that a conference will be held on
	// the 12th and 20th of each month. The value range is [1,31]. If there is no such value in the current month, then
	// for the end of the month.
	Points []int `json:"point,omitempty"`
	// Support the user to specify the number of days N for advance conference notice, the booker will receive the notice
	// of the whole cycle conference, and all the participants will receive the conference notice (including the calendar)
	// N days before each sub-conference time.
	// The input of the number of days N is automatically adjusted according to the interval.
	// If it is held every 2 days on a daily basis, N will automatically become 2, and if it is a Monday or Tuesday
	// every 2 weeks on a weekly basis, N will automatically become 14. Constraints: DST handling is not considered for
	// now. The valid value is range from 0 to 30. The default is 1.
	PreRemindDays *int `json:"preRemindDays,omitempty"`
}

type Configuration struct {
	// Whether to send conference email notification. The default value is determined by the enterprise-level configuration.
	//   true: required.
	//   false: not required.
	IsSendNotify *bool `json:"isSendNotify,omitempty"`
	// Whether to send conference SMS notification. The default value is determined by the enterprise-level configuration.
	//   true: required.
	//   false: not required.
	// Only official commercial enterprises have the right to send conference SMS notifications.
	// Free enterprises will not send conference SMS notifications even if isSendSms is set to true.
	IsSendSms *bool `json:"isSendSms,omitempty"`
	// Whether to send conference calendar notifications. The default value is determined by the enterprise-level configuration.
	//   true: required.
	//   false: not required.
	IsSendCalendar *bool `json:"isSendCalendar,omitempty"`
	// Whether the soft terminal is automatically muted when the guest joins the conference.
	// The default value is determined by the enterprise-level configuration.
	//   true: Automatic mute.
	//   false: Do not mute automatically.
	IsAutoMute *bool `json:"isAutoMute,omitempty"`
	// Whether the guest joins the conference, whether the hard terminal is automatically muted.
	// The default value is determined by the enterprise-level configuration.
	//   true: Automatic mute.
	//   false: Do not mute automatically.
	IsHardTerminalAutoMute *bool `json:"isHardTerminalAutoMute,omitempty"`
	// Whether the guest is password-free (only valid for random conferences).
	//   true: no password.
	//   false: A password is required.
	IsGuestFreePwd *bool `json:"isGuestFreePwd,omitempty"`
	// The range to allow incoming calls.
	//   0: All users.
	//   2: Users within the enterprise.
	//   3: The invited user.
	CallInRestriction *int `json:"callInRestriction,omitempty"`
	// Whether to allow guests to start conferences (only valid for random ID conferences).
	//   true: Allows guests to start conferences.
	//   false: Disables guests from starting conferences.
	AllowGuestStartConf *bool `json:"allowGuestStartConf,omitempty"`
	// Guest password (pure number 4-16 digits long).
	GuestPwd string `json:"guestPwd,omitempty"`
	// Cloud conference room conference ID mode.
	//   0: Fixed ID.
	//   1: Random ID.
	VmrIDType *int `json:"vmrIDType,omitempty"`
	// Automatically extend the conference duration (recommended value range is 0-60).
	//   0: Indicates that the conference ends automatically at the end of the session, and does not extend the conference.
	// Others: Indicates the duration of the automatic extension.
	// Automatically ending the conference is calculated according to the duration of the conference. For example, a scheduled conference starts at 9:00 and ends at 11:00, and the conference lasts for 2 hours. If a participant joins the conference at 8:00, the conference will automatically end at 10:00.
	ProlongLength *int `json:"prolongLength,omitempty"`
	// Whether to open the waiting room (only valid for RTC enterprises).
	//   true: On.
	//   false: not enabled.
	EnableWaitingRoom *bool `json:"enableWaitingRoom,omitempty"`
}

// Create is a method to initiate a conference using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) ([]Conference, error) {
	url := rootURL(c)
	// Whether the conference is the the cycle conference.
	if opts.CycleParams != nil {
		url = cycleURL(c)
	}

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r []Conference
	_, err = c.Post(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return r, err
}

type GetOpts struct {
	// Conference ID.
	ConferenceId string `q:"conferenceID"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// User ID.
	UserId string `q:"userUUID"`
	// Authorization.
	Token string `json:"-"`
}

// Get is a method to obtain scheduled or online conference details using given parameters.
func Get(c *golangsdk.ServiceClient, opts GetOpts) (*GetResp, error) {
	url := showURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r GetResp
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return &r, err
}

// GetOnline is a method to obtain online conference using given parameters.
func GetOnline(c *golangsdk.ServiceClient, opts GetOpts) (*GetResp, error) {
	url := onlineURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r GetResp
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return &r, err
}

type GetHistoryOpts struct {
	// Conference UUID.
	ConferenceUuid string `q:"confUUID"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// User ID.
	UserId string `q:"userUUID"`
	// Authorization token.
	Token string `json:"-"`
}

// GetHistory is a method to obtain history conference using given parameters.
func GetHistory(c *golangsdk.ServiceClient, opts GetHistoryOpts) (*GetResp, error) {
	url := historyURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r GetResp
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return &r, err
}

type ListHistoryOpts struct {
	// User UUID.
	UserUUID string `q:"userUUID"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// Query historical conferences based on the conference subject, the reservation and the string of conference ID
	// keywords.
	SearchKey string `q:"searchKey"`
	// Whether to query the meeting records of all users under the enterprise.
	// If the login account is not an enterprise administrator, this field is invalid.
	// If this field is true, the userUUID field has no effect.
	// The default value is 'false'.
	QueryAll bool `q:"queryAll"`
	// The query's starting date in milliseconds. For example: 1583078400000
	StartDate int `q:"startDate"`
	// The query deadline in milliseconds. For example: 1585756799000
	EndDate int `q:"endDate"`
	//ASC_StartTIME: Sort by meeting start time in ascending order.
	// DSC_StartTIME: Sort in descending order according to the conference start time.
	// ASC_RecordTYPE: Sort according to whether there are recording files or not, and then sort according to the
	//                 conference start time in ascending order by default.
	// DSC_RecordTYPE: Sort according to whether there are recording files or not, and then sort by the conference
	//                 start time in descending order by default.
	SortType string `q:"sortType"`
	// Authorization token.
	Token string `json:"-"`
}

// ListHistory is a method to obtain history conference list using given parameters.
func ListHistory(c *golangsdk.ServiceClient, opts ListHistoryOpts) ([]Conference, error) {
	url := historiesURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r struct {
		Conferences []Conference `json:"data"`
	}
	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return r.Conferences, err
}

type UpdateOpts struct {
	// Conference ID.
	ConferenceID string `q:"conferenceID"`
	// Conference UUID.
	UserUUID string `q:"userUUID"`
	// The conference start time (UTC time).
	//   When creating a reservation conference, if the start time is not specified or the blank string is filled in, it
	//   means that the conference will start immediately.
	//   The time is in UTC, which is the time in time zone 0.
	// Format: yyyy-MM-dd HH:mm
	StartTime string `json:"startTime,omitempty"`
	// The duration of the conference, in minutes, with a maximum value of 1440 and a minimum value of 15.
	// The default is 30.
	Length int `json:"length,omitempty"`
	// Conference subject. The length is limited to 128 characters.
	Subject string `json:"subject,omitempty"`
	// The media type of the conference.
	// It consists of one or more enumeration Strings. When there are multiple enumerations, each enumeration value is
	// separated by "," commas. The enumeration values are as follows:
	//   "Voice": Voice.
	//   "Video": SD video.
	//   "HDVideo": High-definition video (mutually exclusive with Video, if Video and HDVideo are selected at the same
	//     time, the system will select Video by default).
	//   "Telepresence": Telepresence (mutually exclusive with HDVideo and Video, if selected at the same time, the
	//     system uses Telepresence). (reserved field)
	//   "Data": Multimedia (system configuration determines whether to automatically add Data).
	MediaTypes string `json:"mediaTypes,omitempty"`
	// When the soft terminal creates an instant conference, it carries a temporary group ID in the current field, which
	// is carried by the server in the conference-info header field when inviting other participants.
	// The length is limited to 31 characters.
	GroupUri string `json:"groupuri,omitempty"`
	// Attendee list. This list can be used to send conference notifications, conference reminders, and automatic invitations
	// when the conference starts.
	Participants []Participant `json:"attendees,omitempty"`
	// The parameter of the periodic conference. When the conference is a periodic conference, this parameter must be filled in,
	// otherwise the server ignores this parameter.
	// This parameter includes the start date, end date of the periodic conference, the period of the conference and the
	// conference time point in the period.
	CycleParams *CycleParams `json:"cycleParams,omitempty"`
	// Whether the conference automatically starts recording, it only takes effect when the recording type is:
	//   1: Automatically start recording.
	//   0: Do not start recording automatically.
	// The default is not to start automatically.
	IsAutoRecord *int `json:"isAutoRecord,omitempty"`
	// Conference media encryption mode.
	//   0: Adaptive encryption.
	//   1 : Force encryption.
	//   2 : Do not encrypt.
	// Default values are populated by enterprise-level configuration.
	EncryptMode *int `json:"encryptMode,omitempty"`
	// The default language of the conference, the default value is defined by the conference cloud service.
	// For languages supported by the system, it is passed according to the RFC3066 specification.
	//   zh-CN: Simplified Chinese.
	//   en-US: US English.
	Language string `json:"language,omitempty"`
	// The time zone information of the conference time in the conference notification.
	// For time zone information, refer to the time zone mapping relationship.
	// For example: "timeZoneID":"26", the time in the conference notification sent through HUAWEI CLOUD conference will be
	// marked as "2021/11/11 Thursday 00:00 - 02:00 (GMT) Greenwich Standard When: Dublin, Edinburgh, Lisbon, London".
	// For an aperiodic conference, if the conference notification is sent through a third-party system, this field does not
	// need to be filled in.
	TimeZoneID string `json:"timeZoneID,omitempty"`
	// Recording type.
	//   0: Disabled.
	//   1: Live broadcast.
	//   2: Record and broadcast.
	//   3: Live + Recording.
	// Default is disabled.
	RecordType *int `json:"recordType,omitempty"`
	// The mainstream live broadcast address, with a maximum of 255 characters.
	// It is valid when the recording type is 2 or 3.
	LiveAddress string `json:"liveAddress,omitempty"`
	// Auxiliary streaming address, the maximum length is 255 characters.
	// It is valid when the recording type is 2 or 3.
	AuxAddress string `json:"auxAddress,omitempty"`
	// Whether to record auxiliary stream.
	//   0: Do not record.
	//   1: Record.
	// It is valid when the recording type is 2 or 3.
	RecordAuxStream *int `json:"recordAuxStream,omitempty"`
	// Other configuration information for the conference.
	Configuration *Configuration `json:"confConfigInfo,omitempty"`
	// Recording authentication method.
	//   0: Viewable/downloadable via link.
	//   1: Enterprise users can watch/download.
	//   2: Attendees can watch/download.
	// It is valid when the recording type is 2 or 3.
	RecordAuthType *int `json:"recordAuthType,omitempty"`
	// Whether to use the cloud conference room or personal conference ID to hold a reservation conference.
	//   0: Do not use cloud conference room.
	//   1: Use cloud conference room or personal conference ID.
	// Cloud conference rooms are not used by default.
	VmrFlag *int `json:"vmrFlag,omitempty"`
	// VMR ID bound to the current founding account.
	// Obtained by querying the cloud conference room and personal conference ID interface.
	// Note:
	//   vmrID takes the ID returned in the above query interface, not vmrId.
	//   When creating a conference with a personal conference ID, use VMR with vmrMode=0;
	//   When creating a conference in a cloud conference room, use VMR with vmrMode=1.
	VmrId string `json:"vmrID,omitempty"`
	// The number of parties in the conference, the maximum number of participants in the conference.
	//   0: Unlimited.
	// Greater than 0: the maximum number of participants in the conference.
	ConcurrentParticipants *int `json:"concurrentParticipants,omitempty"`
	// The authorization token.
	Token string `json:"-" required:"true"`
}

// Update is a method to udpate the conference configuration using given parameters.
func Update(c *golangsdk.ServiceClient, opts UpdateOpts) ([]Conference, error) {
	url := rootURL(c)
	// Whether the conference is the the cycle conference.
	if opts.CycleParams != nil {
		url = cycleURL(c)
	}
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r []Conference
	_, err = c.Put(url, b, &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return r, err
}

type DeleteOpts struct {
	ConferenceId string `q:"conferenceID"`
	UserId       string `q:"userUUID"`
	Type         int    `q:"type"`
	IsCycle      bool   `json:"-"`
	Token        string `json:"-"`
}

// Delete is a method to cancel conference using given parameters.
func Delete(c *golangsdk.ServiceClient, opts DeleteOpts) error {
	url := rootURL(c)
	// Whether the conference is the the cycle conference.
	if opts.IsCycle {
		url = cycleURL(c)
	}
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()

	_, err = c.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":   "application/json;charset=UTF-8",
			"X-Access-Token": opts.Token,
		},
	})
	return err
}

type AuthOpts struct {
	ConferenceId string `q:"conferenceID"`
	Password     string `json:"-"`
}

// GetControlToken is a method to generate a control token according to the chair password and conference ID.
func GetControlToken(c *golangsdk.ServiceClient, opts AuthOpts) (*AuthResp, error) {
	url := controlURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r AuthResp
	_, err = c.Get(rootURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"X-Password":   opts.Password,
			"X-Login-Type": "1",
		},
	})
	return &r, err
}

type ExtendOpts struct {
	ConferenceId string `q:"conferenceID"`
	IsAuto       int    `json:"auto" required:"true"`
	Duration     int    `json:"duration,omitempty"`
	Token        string `json:"-"`
}

// ExtendTime is a method to extend conference time using given parameters.
func ExtendTime(c *golangsdk.ServiceClient, opts ExtendOpts) error {
	url := controlURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(url, b, nil, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":               "application/json;charset=UTF-8",
			"X-Conference-Authorization": opts.Token,
		},
	})
	return err
}
