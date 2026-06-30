package wizard

import ssh "github.com/duong6003/ssh-wizard/internal/ssh"

type AuthMethod string

const (
	AuthExisting AuthMethod = "existing"
	AuthGenerate AuthMethod = "generate"
	AuthPassword AuthMethod = "password"
)

type ServerConfig struct {
	Alias    string
	Hostname string
	Username string
	Port     int
}

type KeyConfig struct {
	Method         AuthMethod
	KeyType        ssh.KeyType
	PrivateKeyPath string
	PublicKeyPath  string
	Fingerprint    string
	HasPassphrase  bool
}

type State struct {
	Server            *ServerConfig
	Key               *KeyConfig
	KeyInstalled      bool
	ConfigWritten     bool
	ConnectionSuccess bool
}

func NewState() *State {
	return &State{}
}
