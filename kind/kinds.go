package kind

const (
	ProfileMetadata          = 0
	TextNote                 = 1
	RecommendServer          = 2
	FollowList               = 3
	EncryptedDirectMessage   = 4
	Deletion                 = 5
	Repost                   = 6
	Reaction                 = 7
	BadgeAward               = 8
	SimpleGroupChatMessage   = 9
	SimpleGroupThreadedReply = 10
	SimpleGroupThread        = 11
	SimpleGroupReply         = 12
	Seal                     = 13
	DirectMessage            = 14
	GenericRepost            = 16
	ReactionToWebsite        = 17
	ChannelCreation          = 40
	ChannelMetadata          = 41
	ChannelMessage           = 42
	ChannelHideMessage       = 43
	ChannelMuteUser          = 44
	Chess                    = 64
	MergeRequests            = 818
	Bid                      = 1021
	BidConfirmation          = 1022
	OpenTimestamps           = 1040
	GiftWrap                 = 1059
	FileMetadata             = 1063
	LiveChatMessage          = 1311
	Patch                    = 1617
	Issue                    = 1621
	Reply                    = 1622
	StatusOpen               = 1630
	StatusApplied            = 1631
	StatusClosed             = 1632
	StatusDraft              = 1633
	ProblemTracker           = 1971
	Reporting                = 1984
	Label                    = 1985
	RelayReviews             = 1986
	AIEmbeddings             = 1987
	Torrent                  = 2003
	TorrentComment           = 2004
	CoinjoinPool             = 2022
	CommunityPostApproval    = 4550
	JobFeedback              = 7000
	SimpleGroupPutUser       = 9000
	SimpleGroupRemoveUser    = 9001
	SimpleGroupEditMetadata  = 9002
	SimpleGroupDeleteEvent   = 9005
	SimpleGroupCreateGroup   = 9007
	SimpleGroupDeleteGroup   = 9008
	SimpleGroupCreateInvite  = 9009
	SimpleGroupJoinRequest   = 9021
	SimpleGroupLeaveRequest  = 9022
	ZapGoal                  = 9041
	TidalLogin               = 9467
	ZapRequest               = 9734
	Zap                      = 9735
	Highlights               = 9802
	MuteList                 = 10000
	PinList                  = 10001
	RelayListMetadata        = 10002
	BookmarkList             = 10003
	CommunityList            = 10004
	PublicChatList           = 10005
	BlockedRelayList         = 10006
	SearchRelayList          = 10007
	SimpleGroupList          = 10009
	InterestList             = 10015
	EmojiList                = 10030
	DMRelayList              = 10050
	UserServerList           = 10063
	FileStorageServerList    = 10096
	GoodWikiAuthorList       = 10101
	GoodWikiRelayList        = 10102
	NWCWalletInfo            = 13194
	LightningPubRPC          = 21000
	ClientAuthentication     = 22242
	NWCWalletRequest         = 23194
	NWCWalletResponse        = 23195
	NostrConnect             = 24133
	Blobs                    = 24242
	HTTPAuth                 = 27235
	CategorizedPeopleList    = 30000
	CategorizedBookmarksList = 30001
	RelaySets                = 30002
	BookmarkSets             = 30003
	CuratedSets              = 30004
	CuratedVideoSets         = 30005
	MuteSets                 = 30007
	ProfileBadges            = 30008
	BadgeDefinition          = 30009
	InterestSets             = 30015
	StallDefinition          = 30017
	ProductDefinition        = 30018
	MarketplaceUI            = 30019
	ProductSoldAsAuction     = 30020
	Article                  = 30023
	DraftArticle             = 30024
	EmojiSets                = 30030
	ModularArticleHeader     = 30040
	ModularArticleContent    = 30041
	ReleaseArtifactSets      = 30063
	ApplicationSpecificData  = 30078
	LiveEvent                = 30311
	UserStatuses             = 30315
	ClassifiedListing        = 30402
	DraftClassifiedListing   = 30403
	RepositoryAnnouncement   = 30617
	RepositoryState          = 30618
	SimpleGroupMetadata      = 39000
	SimpleGroupAdmins        = 39001
	SimpleGroupMembers       = 39002
	SimpleGroupRoles         = 39003
	WikiArticle              = 30818
	Redirects                = 30819
	Feed                     = 31890
	DateCalendarEvent        = 31922
	TimeCalendarEvent        = 31923
	Calendar                 = 31924
	CalendarEventRSVP        = 31925
	HandlerRecommendation    = 31989
	HandlerInformation       = 31990
	VideoEvent               = 34235
	ShortVideoEvent          = 34236
	VideoViewEvent           = 34237
	CommunityDefinition      = 34550
)

func IsRegularKind(kind int) bool {
	return kind < 10000 && kind != 0 && kind != 3
}

func IsReplaceableKind(kind int) bool {
	return kind == 0 || kind == 3 || (10000 <= kind && kind < 20000)
}

func IsEphemeralKind(kind int) bool {
	return 20000 <= kind && kind < 30000
}

func IsAddressableKind(kind int) bool {
	return 30000 <= kind && kind < 40000
}

var Text = []int{
	ProfileMetadata,
	TextNote,
	Article,
	SimpleGroupThread,
	Reply,
	Repost,
	Issue,
	Reply,
	MergeRequests,
	WikiArticle,
	Issue,
	StatusOpen,
	StatusApplied,
	StatusClosed,
	StatusDraft,
	Torrent,
	TorrentComment,
	DateCalendarEvent,
	TimeCalendarEvent,
	Calendar,
	CalendarEventRSVP,
}

func IsText(ki int) bool {
	for _, v := range Text {
		if v == ki {
			return true
		}
	}
	return false
}
