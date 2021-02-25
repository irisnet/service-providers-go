package service

// AddKey adds a key with the specified name and password
func (s ServiceClientWrapper) AddKey(name string, password string) (addr string, mnemonic string, err error) {
	return s.ServiceClient.Insert(name, password)
}

// ShowKey queries the given key
func (s ServiceClientWrapper) ShowKey(name string, password string) (addr string, err error) {
	_, address, err := s.ServiceClient.Find(name, password)
	return address.String(), err
}

// ImportKey imports the specified key
func (s ServiceClientWrapper) ImportKey(name string, password string, keyArmor string) (addr string, err error) {
	return s.ServiceClient.Import(name, password, keyArmor)
}
