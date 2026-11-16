package validations

type Login struct { // login validation rules
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePassword struct { // change password validation rules
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPassword struct { // reset password validation rules
	Email string `json:"email" binding:"required,email"`
}

type VerifyPassword struct { // verify reset password validation rules
	Token    string `json:"token" binding:"required,min=32,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}
