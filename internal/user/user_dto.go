package user

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (UserResponse) FromEntity(user UserEntity) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
