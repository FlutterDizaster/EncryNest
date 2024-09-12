package keychain

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

type MasterKey []byte
