package identity

type Identity struct {
	noxId    string
	name     string
	deviceId string
}

func BotIdentity(noxId string, name string, deviceId string) Identity {
	return Identity{
		noxId:    noxId,
		name:     name,
		deviceId: deviceId,
	}
}

func (i *Identity) NoxId() string {
	return i.noxId
}

func (i *Identity) Name() string {
	return i.name
}

func (i *Identity) DeviceId() string {
	return i.deviceId
}
