package iface

type Connectioner interface {
	SendMsg(data []byte) error
}
