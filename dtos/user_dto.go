package dto

type CreateUser struct {
	UserName     string `json:"username"`
	Name         string `json:"name"`
	TempPassword string `json:"tempPassword"` // hashed password from fe
}

type Register struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	EccKey   string `json:"eccKey"`
	RsaKey   string `json:"rsaKey"`
}

type LoginReturn struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

type CreateChallenge struct {
	PublicKey string `json:"publicKey"`
}

type VerifyChallenge struct {
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}
