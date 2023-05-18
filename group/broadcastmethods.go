package group

func (brdreply *BroadCastReplyMessage) ADDToPeerTable() {
	groupPeer := &GroupPeer{
		PeerId:   brdreply.PeerId,
		UserName: "Test",
	}
	PeerTable = append(PeerTable, *groupPeer)
}
