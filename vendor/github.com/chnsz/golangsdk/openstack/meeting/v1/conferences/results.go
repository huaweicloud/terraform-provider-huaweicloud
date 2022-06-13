package conferences

type Conference struct {
	// Meeting ID. The length is limited to 32 characters.
	ID string `json:"conferenceID"`
	// Conference theme. The length is limited to 128 characters.
	Subject string `json:"subject"`
	// The number of parties in the meeting.
	Size int `json:"size"`
	// The time zone information of the meeting time in the meeting notification.
	// For time zone information, refer to the time zone mapping relationship.
	// For example: "timeZoneID": "26", the time in the meeting notice sent through HUAWEI CLOUD meeting will be marked
	// as "2021/11/11 Thursday 00:00 - 02:00 (GMT) Greenwich Standard When: Dublin, Edinburgh, Lisbon, London".
	TimeZoneId string `json:"timeZoneID"`
	// The meeting start time (YYYY-MM-DD HH:MM).
	StartTime string `json:"startTime"`
	// The meeting end time (YYYY-MM-DD HH:MM).
	EndTime string `json:"endTime"`
	// The media type of the meeting.
	// It consists of one or more enumeration Strings.
	// When there are multiple enumerations, each enumeration value is separated by "," commas.
	//   Voice: Voice.
	//   Video: SD video.
	//   HDVideo: high-definition video (mutually exclusive with "Video", if "Video" and "HDVideo" are selected at the
	//            same time, the system selects "Video" by default).
	//   Telepresence: Telepresence (mutually exclusive with "HDVideo" and "Video", if selected at the same time, the
	//                 system uses "Telepresence"). (reserved field)
	//   Data: Multimedia.
	MediaTypes string `json:"mediaTypes"`
	// Currently, only the Created and Scheduled states will be returned.  If the meeting has been held, the Created
	// state will be returned, otherwise, the Scheduled state will be returned.
	//   Schedule: Schedule status.
	//   Creating: The state is being created.
	//   Created: The meeting has been created and is in progress.
	//   Destroyed: The meeting has been closed.
	ConferenceState string `json:"conferenceState"`
	// Conference language.
	Language string `json:"language"`
	// Conference access code.
	AccessNumber string `json:"accessNumber"`
	// Meeting password entry. The subscriber returns the host password and guest password.
	// The host password is returned when the host queries.
	// The guest password is returned when the guest is queried.
	PasswordEntry []PasswordEntry `json:"passwordEntry"`
	// The UUID of the meeting booker.
	UserUUID string `json:"userUUID"`
	// Meeting booker account name. The maximum length is limited to 96 characters.
	ScheduserName string `json:"scheduserName"`
	// Conference type (Parameters used by the front-end).
	//   0 : Normal meeting.
	//   2 : Periodic meeting.
	ConferenceType int `json:"conferenceType"`
	// Conference type.
	//   FUTURE
	//   IMMEDIATELY
	//   CYCLE
	ConfType string `json:"confType"`
	// Periodic meeting parameters. Carry this parameter when the conference is a periodic conference.
	// This parameter includes the start date, end date of the periodic meeting, the period of the meeting and the
	// meeting time point in the period.
	CycleParams CycleParams `json:"cycleParams"`
	// Whether to automatically mute the session.
	//   0 : Do not mute automatically.
	//   1 : Auto mute.
	IsAutoMute int `json:"isAutoMute"`
	// Whether to automatically start recording.
	//   0 : Do not start automatically.
	//   1 : start automatically.
	IsAutoRecord int `json:"isAutoRecord"`
	// Host meeting link address.
	ChairJoinUri string `json:"chairJoinUri"`
	// Common attendee meeting link address. The maximum length is 1024.
	GuestJoinUri string `json:"guestJoinUri"`
	// Audience meeting link address. The maximum length is 1024. (webinar scenario)
	AudienceJoinUri string `json:"audienceJoinUri"`
	// Recording type.
	//   0: Disabled.
	//   1: Live broadcast.
	//   2: Record and broadcast.
	//   3: Live + Recording.
	RecordType int `json:"recordType"`
	// Auxiliary stream live address.
	AuxAddress string `json:"auxAddress"`
	// Mainstream live broadcast address.
	LiveAddress string `json:"liveAddress"`
	// Whether to record auxiliary streams.
	//   0: No.
	//   1: Yes.
	RecordAuxStream int `json:"recordAuxStream"`
	// Recording and broadcasting authentication method.
	// The recording type is: recording, live+recording, and it is valid.
	//   0: Viewable/downloadable via link.
	//   1: Enterprise users can watch/download.
	//   2: Attendees can watch/download.
	RecordAuthType int `json:"recordAuthType"`
	// Live address. (It will be returned when the live room is configured)
	LiveUrl string `json:"liveUrl"`
	// Other configuration information for the conference.
	Configuration Configuration `json:"confConfigInfo"`
	// Whether to use the cloud conference room to hold a reservation meeting.
	//   0: Do not use cloud conference room.
	//   1: Use cloud conference room.
	// The interface shows that the conference ID needs to use "vmrConferenceID" as the conference ID;
	// the "conferenceID" field is still used for conference business operations such as querying conference details,
	// logging in to conference control, and typing in a conference.
	VmrFlag int `json:"vmrFlag"`
	// Only the historical conference return value is valid. There are no recording files by default.
	//   true: There is a recording file.
	//   false: No recording file.
	IsHasRecordFile bool `json:"isHasRecordFile"`
	// The conference ID of the cloud conference room. If "vmrFlag" is "1", this field is not empty.
	VmrConferenceId string `json:"vmrConferenceID"`
	// The UUID of the meeting.
	// The UUID is only returned when a meeting that starts immediately is created.
	// If it is a future meeting, the UUID will not be returned.
	// You can get the UUID of the historical conference through "Query Historical Conference List".
	ConfUUID string `json:"confUUID"`
	// Information about some of the invited participants.
	// Only the first 20 soft terminal participant information and the first 20 hard terminal participant information
	// are returned. Do not return the information of participants who actively joined in the conference.
	// The "Query Conference List" and "Query Conference Details" interfaces return the participants who were invited
	// when the conference was scheduled and the participants who were invited by the host in the conference.
	// The "Query Online Conference List", "Query Online Conference Details", "Query History Conference List" and "Query
	// History Conference Details" interfaces return the participants who were invited when the conference was
	// scheduled. Attendees invited by the host in the meeting are not returned.
	Participants []ParticipantResp `json:"partAttendeeInfo"`
	// Number of hard terminals, such as IdeaHub, TE30, etc.
	TerminlCount int `json:"terminlCount"`
	// The number of common terminals, such as PC terminal, mobile terminal app, etc.
	NormalCount int `json:"normalCount"`
	// The business name of the meeting booker. Maximum length 96.
	DeptName string `json:"deptName"`
	// Attendee role.
	//   chair : Chair.
	//   general : The guest.
	//   audience : The audience.
	Role string `json:"role"`
	// Identifies whether it is a multi-stream video conference.
	// 1 : Multi-stream conference.
	MultiStreamFlag int `json:"multiStreamFlag"`
	// Webinar or not.
	Webinar bool `json:"webinar"`
	// Type of meeting.
	// COMMON : Normal conference.
	// RTC : RTC conference.
	ConfMode string `json:"confMode"`
	// VMR appointment record.
	// true : VMR conference.
	// false : Normal meeting.
	ScheduleVmr bool `json:"scheduleVmr"`
	// The UUID of the cloud meeting room.
	VmrId string `json:"vmrID"`
	// The number of parties in the conference, the maximum number of participants in the conference.
	ConcurrentParticipants int `json:"concurrentParticipants"`
	// Current multi-screen information.
	PicDisplay MultipicDisplayDo `json:"picDisplay"`
	// List of periodic sub-conferences.
	Subconferences []SubConference `json:"subConfs"`
	// The UUID of the first cycle subconference.
	CycleSubConfId string `json:"cycleSubConfID"`
}

type PasswordEntry struct {
	// Conference role.
	//   chair: The host of the meeting.
	//   general: General participants.
	ConferenceRole string `json:"conferenceRole"`
	// The password (clear text) of the role in the meeting.
	Password string `json:"password"`
}

type ParticipantResp struct {
	// Attendee name or nickname. The length is limited to 96 characters.
	Name string `json:"name"`
	// Phone number (support SIP, TEL number format). Maximum of 127 characters.
	// At least one of phone, email and sms must be filled in.
	// When "type" is "telepresence" and the device is a three-screen telepresence, fill in the number of the middle
	// screen in this field. (Three-screen telepresence is a reserved field)
	Phone string `json:"phone"`
	// The default value is defined by the conference AS, and the number type is enumerated as follows:
	//   "normal": soft terminal.
	//   "telepresence": telepresence. Single-screen and triple-screen telepresence belong to this category. (reserved)
	//   "terminal": conference room or hard terminal.
	//   "outside": The outside participant.
	//   "mobile": User's mobile phone number.
	//   "telephone": The user's landline phone. (reserved field)
	Type string `json:"type"`
}

type MultipicDisplayDo struct {
	// Whether to set multi-screen manually.
	//   0: The system automatically multi-screen.
	//   1: Manually set multi-screen.
	ManualSet int `json:"manualSet"`
	// Screen Type, value range:
	//   Single: single screen
	//   Two: Two pictures
	//   Three: Three pictures
	//     Three-2: Three pictures-2
	//     Three-3: Three pictures-3
	//     Three-4: Three pictures-4
	//   Four: Quad Picture
	//     Four-2: Quad Picture-2
	//     Four-3: Quad Picture-3
	//   Five: Five pictures
	//     Five-2: Five pictures-2
	//   Six: Six pictures
	//     Six-2: Six pictures-2
	//     Six-3: Six pictures-3
	//     Six-4: Six pictures-4
	//     Six-5: Six pictures-5
	//   Seven: Seven pictures
	//     Seven-2: Seven pictures-2
	//     Seven-3: Seven pictures-3
	//     Seven-4: Seven pictures-4
	//   Eight: Eight-picture
	//     Eight-2: Eight-picture-2
	//     Eight-3: Eight-picture-3
	//     Eight-4: Eight-picture-4
	//   Nine: Nine pictures
	//   Ten: Ten pictures
	//     Ten-2: Ten pictures-2
	//     Ten-3: Ten pictures-3
	//     Ten-4: Ten pictures-4
	//     Ten-5: Ten pictures-5
	//     Ten-6: Ten pictures-6
	//   Thirteen: Thirteen pictures
	//     Thirteen-2: Thirteen pictures-2
	//     Thirteen-3: Thirteen pictures-3
	//     Thirteen-4: Thirteen pictures-4
	//     Thirteen-5: Thirteen pictures-5
	//     ThirteenR: Thirteen Frames R
	//     ThirteenM: Thirteen Frames M
	//   Sixteen: Sixteen screen
	//   Seventeen: Seventeen pictures
	//   Twenty-Five: Twenty-Five
	//   Custom: custom multi-screen
	ImageType string `json:"imageType"`
	// Subscreen list.
	SubscriberInPics []PicInfoNotify `json:"subscriberInPics"`
	// Indicates the polling interval, in seconds.
	// This parameter is valid when there are multiple video sources in the same sub-picture.
	SwitchTime string `json:"switchTime"`
	// Customize multi-screen layout information.
	PicLayout PicLayout `json:"picLayoutInfo"`
}

type PicInfoNotify struct {
	// The number of each screen in the multi-screen, starting from 1.
	Index int `json:"index"`
	// Session ID in each screen, namely callNumber.
	Id []string `json:"id"`
	// Whether it is an auxiliary stream.
	//   0: Not an auxiliary stream.
	//   1: It is an auxiliary stream.
	Share int `json:"share"`
}

type PicLayout struct {
	// The number of horizontal small grids.
	X int `json:"x"`
	// The number of vertical small grids.
	Y int `json:"y"`
	// Multi-screen information.
	SubPicLayouts []SubPicLayout `json:"subPicLayoutInfoList"`
}

type SubPicLayout struct {
	// Subpicture index.
	Id int `json:"id"`
	// The index of the sprite from left to right.
	Left int `json:"left"`
	// The index of the sprite from top to bottom.
	Top int `json:"top"`
	// The horizontal size of the sprite.
	XSize int `json:"xSize"`
	// The vertical size of the sprite.
	YSize int `json:"ySize"`
}

type SubConference struct {
	// Subconference UUID
	UUID string `json:"cycleSubConfID"`
	// Conference ID, the length is limited to no more than 32 characters
	ID string `json:"conferenceID"`
	// The media type of the meeting. It consists of one or more enumeration Strings.
	// When there are multiple enumerations, each enumeration value is separated by "," commas.
	// The enumeration values are as follows:
	//   Voice: Voice
	//   Video: SD video
	//   HDVideo: high-definition video (mutually exclusive with Video, if Video and HDVideo are selected at the same
	//            time, the system defaults to Video)
	//   Telepresence: Telepresence (mutually exclusive with HDVideo and Video, if selected at the same time, the system
	//                 uses Telepresence) - not supported yet
	//   Data: Multimedia
	MediaType string `json:"mediaType"`
	// Conference start time (format: YYYY-MM-DD HH:MM)
	StartTime string `json:"startTime"`
	// Conference end time (format: YYYY-MM-DD HH:MM)
	EndTime string `json:"endTime"`
	// Whether to automatically start recording.
	IsAutoRecord int `json:"isAutoRecord"`
	// The recording and broadcasting authentication method is valid when the recording and broadcasting types are:
	// recording and broadcasting, live broadcast + recording and broadcasting.
	//   0 is the old authentication method, and the url carries token authentication.
	//   1 is user authentication for intra-enterprise conferences.
	//   2 is the user authentication for the conference within the conference.
	RecordAuthType int `json:"recordAuthType"`
	// Meeting description, the length is limited to 200 characters.
	Description string `json:"description"`
	// Other configuration information of periodic subconferences.
	SubConfiguration SubConfiguration `json:"confConfigInfo"`
}

type SubConfiguration struct {
	// The range to allow incoming calls.
	//   0: All users.
	//   2: Users within the enterprise.
	//   3: The invited user.
	CallInRestriction int `json:"callInRestriction"`
	// The range that the webinar audience is allowed to call in.
	//   0 for all users.
	//   2 Enterprise users and invited users.
	AudienceCallInRestriction int `json:"audienceCallInRestriction"`
	// Whether to allow guests to start meetings (only random meeting IDs are valid).
	//   true: Allows guests to start meetings.
	//   false: Disables guests from starting meetings.
	AllowGuestStartConf bool `json:"allowGuestStartConf"`
	// Whether to enable waiting room.
	EnableWaitingRoom bool `json:"enableWaitingRoom"`
	// Webinar Audience Display Strategy.
	ShowAudiencePolicy ShowAudiencePolicy `json:"showAudienceCountInfo"`
}

type ShowAudiencePolicy struct {
	// Audience display strategy: The server is used to calculate the number of audiences and send it to the client to
	// control the audience display.
	//   0: do not display.
	//   1: Multiply display the number of participants, based on the real-time number of participants or the cumulative
	// number of participants (assuming N), the multiplication setting can be performed.
	// Notes: Supports setting the multiplier X and the base number Y.
	// After setting, the number of people displayed is: NX+Y.
	//   X supports setting to 1 decimal place. When NX calculates a non-integer, it will be rounded down.
	//   The range of X is 1~10, and the range of Y is 0~10000.
	ShowAudienceMode int `json:"showAudienceMode"`
	// The basic number of people, the range is 0~10000
	BaseAudienceCount int `json:"baseAudienceCount"`
	// Multiplier, the range is 1~10, it can be set to 1 decimal place
	Multiple float64 `json:"multiple"`
}

type GetResp struct {
	Conference Conference      `json:"conferenceData"`
	Data       PageParticipant `json:"data"`
}

type PageParticipant struct {
	// The number of records per page.
	Limit int `json:"limit"`
	// total number of the participants.
	Count int `json:"count"`
	// The offset of the number of records, how many records there are before this page.
	Offset int `json:"offset"`
	// Invited attendee information. It includes the attendees invited when the meeting was scheduled and the attendees
	// invited by the host during the meeting.
	// Do not return the information of participants who actively joined in the conference.
	Participants []ParticipantDetail `json:"data"`
}

type ParticipantDetail struct {
	// The participant's number.
	ID string `json:"participantID"`
	// The name (nickname) of the attendee.
	Name string `json:"name"`
	// role in the meeting.
	//   1: The meeting host.
	//   0: Regular attendees.v
	Role int `json:"role"`
	// user status. Currently fixed return to MEETING.
	State string `json:"state"`
	// Information about the conference room where the terminal is located. (reserved field)
	Address string `json:"address"`
	// The default value is defined by the conference AS.
	//   normal: soft terminal.
	//   telepresence: telepresence. Single-screen and triple-screen telepresence belong to this category. (reserved field)
	//   terminal: conference room or hard terminal.
	//   outside: The outside participant.
	//   mobile: User's mobile phone number.
	//   telephone: The user's landline phone. (reserved field)
	Type string `json:"attendeeType"`
	// The account ID of the participant.
	// In the case of account/password authentication, it indicates the ID of the HUAWEI CLOUD conference account.
	// If it is an APP ID authentication scenario, it indicates the User ID of a third party.
	AccountId string `json:"accountId"`
	// email address. Maximum of 255 characters.
	Email string `json:"email"`
	// Mobile number for SMS notification. Maximum of 127 characters.
	SMS string `json:"sms"`
	// Department name. Maximum of 96 characters.
	DeptName string `json:"deptName"`
	// Subscriber's user UUID.
	UserUUID string `json:"userUUID"`
	// App ID, application identification, an application only needs to be created once, refer to "App ID application".
	AppId string `json:"appId"`
	// Whether to automatically invite this attendee when the meeting starts.
	// The default value is determined by the enterprise-level configuration.
	//   0: do not automatically invite
	//   1: auto-invite
	IsAutoInvite int `json:"isAutoInvite"`
	// Whether to not superimpose the venue name
	//   true: no overlay
	//   false: overlay
	IsNotOverlayPidName bool `json:"isNotOverlayPidName"`
}

type AuthResp struct {
	// Token information.
	TokenDetail TokenDetail `json:"data"`
	// Query the temporary token in the address book.
	AddressToken string `json:"addressToken"`
	// Global external network IP.
	GloablPublicIP string `json:"gloablPublicIP"`
}

type TokenDetail struct {
	// Conference control authentication token.
	Token string `json:"token"`
	//The websocket chain building authentication Token.
	TmpWsToken string `json:"tmpWsToken"`
	// The websocket linking URL.
	WsURL string `json:"wsURL"`
	// The roles in the conference, the enumeration values are as follows.
	//   0 : Conference chairperson.
	//   1 : Regular attendees.
	Role int `json:"role"`
	// Session expiration time. UTC time in milliseconds.
	ExpireTime int `json:"expireTime"`
	// The ID of the conference booker.
	UserID string `json:"userID"`
	// ID of the company to which the conference belongs.
	OrgID string `json:"orgID"`
	// When requested by the terminal, the site ID after the terminal joins the conference is returned.
	ParticipantID string `json:"participantID"`
	// It will control the time when the token expires. (in seconds)
	ConfTokenExpireTime int `json:"confTokenExpireTime"`
	// The current meeting ID of the cloud conference room meeting.
	VmrCurrentConfID string `json:"vmrCurrentConfID"`
	// Websocket message push support type.
	SupportNotifyType []string `json:"supportNotifyType"`
}

type ErrResponse struct {
	// Error code.
	Code string `json:"error_code"`
	// Error message.
	Message string `json:"error_msg"`
}
