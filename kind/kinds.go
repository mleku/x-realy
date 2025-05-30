package kind

const (
	ProfileMetadata          int = 0
	TextNote                 int = 1
	RecommendServer          int = 2
	FollowList               int = 3
	EncryptedDirectMessage   int = 4
	Deletion                 int = 5
	Repost                   int = 6
	Reaction                 int = 7
	BadgeAward               int = 8
	SimpleGroupChatMessage   int = 9
	SimpleGroupThreadedReply int = 10
	SimpleGroupThread        int = 11
	SimpleGroupReply         int = 12
	Seal                     int = 13
	DirectMessage            int = 14
	GenericRepost            int = 16
	ReactionToWebsite        int = 17
	ChannelCreation          int = 40
	ChannelMetadata          int = 41
	ChannelMessage           int = 42
	ChannelHideMessage       int = 43
	ChannelMuteUser          int = 44
	Chess                    int = 64
	MergeRequests            int = 818
	Bid                      int = 1021
	BidConfirmation          int = 1022
	OpenTimestamps           int = 1040
	GiftWrap                 int = 1059
	FileMetadata             int = 1063
	LiveChatMessage          int = 1311
	Patch                    int = 1617
	Issue                    int = 1621
	Reply                    int = 1622
	StatusOpen               int = 1630
	StatusApplied            int = 1631
	StatusClosed             int = 1632
	StatusDraft              int = 1633
	ProblemTracker           int = 1971
	Reporting                int = 1984
	Label                    int = 1985
	RelayReviews             int = 1986
	AIEmbeddings             int = 1987
	Torrent                  int = 2003
	TorrentComment           int = 2004
	CoinjoinPool             int = 2022
	CommunityPostApproval    int = 4550
	JobFeedback              int = 7000
	SimpleGroupPutUser       int = 9000
	SimpleGroupRemoveUser    int = 9001
	SimpleGroupEditMetadata  int = 9002
	SimpleGroupDeleteEvent   int = 9005
	SimpleGroupCreateGroup   int = 9007
	SimpleGroupDeleteGroup   int = 9008
	SimpleGroupCreateInvite  int = 9009
	SimpleGroupJoinRequest   int = 9021
	SimpleGroupLeaveRequest  int = 9022
	ZapGoal                  int = 9041
	TidalLogin               int = 9467
	ZapRequest               int = 9734
	Zap                      int = 9735
	Highlights               int = 9802
	MuteList                 int = 10000
	PinList                  int = 10001
	RelayListMetadata        int = 10002
	BookmarkList             int = 10003
	CommunityList            int = 10004
	PublicChatList           int = 10005
	BlockedRelayList         int = 10006
	SearchRelayList          int = 10007
	SimpleGroupList          int = 10009
	InterestList             int = 10015
	EmojiList                int = 10030
	DMRelayList              int = 10050
	UserServerList           int = 10063
	FileStorageServerList    int = 10096
	GoodWikiAuthorList       int = 10101
	GoodWikiRelayList        int = 10102
	NWCWalletInfo            int = 13194
	LightningPubRPC          int = 21000
	ClientAuthentication     int = 22242
	NWCWalletRequest         int = 23194
	NWCWalletResponse        int = 23195
	NostrConnect             int = 24133
	Blobs                    int = 24242
	HTTPAuth                 int = 27235
	CategorizedPeopleList    int = 30000
	CategorizedBookmarksList int = 30001
	RelaySets                int = 30002
	BookmarkSets             int = 30003
	CuratedSets              int = 30004
	CuratedVideoSets         int = 30005
	MuteSets                 int = 30007
	ProfileBadges            int = 30008
	BadgeDefinition          int = 30009
	InterestSets             int = 30015
	StallDefinition          int = 30017
	ProductDefinition        int = 30018
	MarketplaceUI            int = 30019
	ProductSoldAsAuction     int = 30020
	Article                  int = 30023
	DraftArticle             int = 30024
	EmojiSets                int = 30030
	ModularArticleHeader     int = 30040
	ModularArticleContent    int = 30041
	ReleaseArtifactSets      int = 30063
	ApplicationSpecificData  int = 30078
	LiveEvent                int = 30311
	UserStatuses             int = 30315
	ClassifiedListing        int = 30402
	DraftClassifiedListing   int = 30403
	RepositoryAnnouncement   int = 30617
	RepositoryState          int = 30618
	SimpleGroupMetadata      int = 39000
	SimpleGroupAdmins        int = 39001
	SimpleGroupMembers       int = 39002
	SimpleGroupRoles         int = 39003
	WikiArticle              int = 30818
	Redirects                int = 30819
	Feed                     int = 31890
	DateCalendarEvent        int = 31922
	TimeCalendarEvent        int = 31923
	Calendar                 int = 31924
	CalendarEventRSVP        int = 31925
	HandlerRecommendation    int = 31989
	HandlerInformation       int = 31990
	VideoEvent               int = 34235
	ShortVideoEvent          int = 34236
	VideoViewEvent           int = 34237
	CommunityDefinition      int = 34550
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
