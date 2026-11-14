package validations

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyPassword struct {
	Token    string `json:"token" binding:"required,min=32,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

type FcmModel struct {
	Title   string   `json:"title" binding:"required"`
	Msg     string   `json:"msg" binding:"required"`
	UserIDs []string `json:"uid" binding:"required"`
}
