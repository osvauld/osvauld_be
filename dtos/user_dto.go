package dto

type CreateUser struct {
	UserName     string `json:"username" binding:"required"`
	Name         string `json:"name" binding:"required"`
	TempPassword string `json:"tempPassword" binding:"required"` // hashed password from fe
}

type Register struct {
	UserName      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	DeviceKey     string `json:"eccKey" binding:"required"`
	EncryptionKey string `json:"rsaKey" binding:"required"`
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
