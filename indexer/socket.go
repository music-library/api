package indexer

type Socket struct {
	Sessions map[string]*SocketSession
}

type SocketSession struct {
	Id             string
	PlayingTrackId string
}

func NewSocket() *Socket {
	return &Socket{Sessions: make(map[string]*SocketSession)}
}

func (s *Socket) GetOrCreateSession(userId string) *SocketSession {
	if _, ok := s.Sessions[userId]; !ok {
		s.Sessions[userId] = &SocketSession{Id: userId}
	}

	return s.Sessions[userId]
}

func (s *Socket) PlayingTracks() []string {
	playingTracks := []string{}
	cacheMap := make(map[string]bool)

	for _, session := range s.Sessions {
		if session.PlayingTrackId != "" {
			if _, ok := cacheMap[session.PlayingTrackId]; !ok {
				cacheMap[session.PlayingTrackId] = true
				playingTracks = append(playingTracks, session.PlayingTrackId)
			}
		}
	}

	return playingTracks
}
