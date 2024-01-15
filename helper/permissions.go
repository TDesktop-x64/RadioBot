package helper

import tdlib "github.com/c0re100/gotdlib/client"

func NewChatMemeberStatus(chatId int64, memberId tdlib.MessageSender, status tdlib.ChatMemberStatus) *tdlib.SetChatMemberStatusRequest {
	return &tdlib.SetChatMemberStatusRequest{
		ChatId:   chatId,
		MemberId: memberId,
		Status:   status,
	}
}

func NewChatMemberStatusMember() *tdlib.ChatMemberStatusMember {
	return &tdlib.ChatMemberStatusMember{}
}

func NewChatMemberStatusBanned(bannedUntilDate int32) *tdlib.ChatMemberStatusBanned {
	return &tdlib.ChatMemberStatusBanned{
		BannedUntilDate: bannedUntilDate,
	}
}

func NewChatPermissions(canSendBasicMessages, canSendAudios, canSendDocuments, canSendPhotos, canSendVideos, canSendVideoNotes, canSendVoiceNotes, canSendPolls, canSendStickers, canSendAnimations, canSendGames, canUseInlineBots, canAddWebPagePreviews, canChangeInfo, canInviteUsers, canPinMessages bool) *tdlib.ChatPermissions {
	return &tdlib.ChatPermissions{
		CanSendBasicMessages:  canSendBasicMessages,
		CanSendAudios:         canSendAudios,
		CanSendDocuments:      canSendDocuments,
		CanSendPhotos:         canSendPhotos,
		CanSendVideos:         canSendVideos,
		CanSendVideoNotes:     canSendVideoNotes,
		CanSendVoiceNotes:     canSendVoiceNotes,
		CanSendPolls:          canSendPolls,
		CanSendStickers:       canSendStickers,
		CanSendAnimations:     canSendAnimations,
		CanSendGames:          canSendGames,
		CanUseInlineBots:      canUseInlineBots,
		CanAddWebPagePreviews: canAddWebPagePreviews,
		CanChangeInfo:         canChangeInfo,
		CanInviteUsers:        canInviteUsers,
		CanPinMessages:        canPinMessages,
	}
}

func NewChatMemberStatusRestricted(isMember bool, restrictedUntilDate int32, permissions *tdlib.ChatPermissions) *tdlib.ChatMemberStatusRestricted {
	return &tdlib.ChatMemberStatusRestricted{
		IsMember:            isMember,
		RestrictedUntilDate: restrictedUntilDate,
		Permissions:         permissions,
	}
}

func NewChatMemberStatusRestrictedAll(isMember bool, restrictedUntilDate int32) *tdlib.ChatMemberStatusRestricted {
	return &tdlib.ChatMemberStatusRestricted{
		IsMember:            isMember,
		RestrictedUntilDate: restrictedUntilDate,
		Permissions:         &tdlib.ChatPermissions{},
	}
}
