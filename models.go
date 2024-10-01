package main

type DeviceInfo struct {
	DevEUI     []byte
	DevAddr    []byte
	NwkSEncKey []byte
}

type RouteDevice struct {
	RouteID    string `json:"route_id"`
	DevAddr    string `json:"devaddr"`
	SessionKey string `json:"session_key"`
	MaxCopies  int    `json:"max_copies"`
}

type DeviceUpdate struct {
	RouteID    string `json:"route_id"`
	DevAddr    string `json:"devaddr"`
	SessionKey string `json:"session_key"`
	MaxCopies  int    `json:"max_copies"`
	Action     string `json:"action"`
}
