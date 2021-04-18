package msg_push

import (
	"sync"
)

// 一个房间代表一个订阅推送类型
type Room struct {
	ID    int
	Title string
	RConn sync.Map
}

func NewRoom(id int, title string) *Room {
	return &Room{
		ID:    id,
		Title: title,
	}
}

// 加入房间
//func (r *Room) JoinRoom(ws *WsConnection) error {
//	if _, ok := r.RConn.Load(ws.GetWsId()); ok{
//		return errors.New("already exists")
//	}
//	r.RConn.Store(ws.GetWsId(), ws)
//	return nil
//}
