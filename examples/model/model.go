package model

type WsUserData struct {
	UserID        int    `json:"user_id"`
	WsProtocalID  int64  `json:"ws_protocal_id"`
	To            int    `json:"to"`
	ProductionExt string `json:"production_ext"`
	Data          map[string]interface{}
}
