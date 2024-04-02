package dto

type CreateUser struct {
	UserName     string `json:"username" binding:"required"`
	Name         string `json:"name" binding:"required"`
	TempPassword string `json:"tempPassword" binding:"required"` // hashed password from fe
	Type         string `json:"type"`
}

type Register struct {
	UserName      string `json:"username" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
	DeviceKey     string `json:"deviceKey" binding:"required"`
	EncryptionKey string `json:"encryptionKey" binding:"required"`
}

type TempLogin struct {
	UserName     string `json:"username" binding:"required"`
	TempPassword string `json:"tempPassword" binding:"required"`
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

type CheckUserAvailability struct {
	UserName string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
}
