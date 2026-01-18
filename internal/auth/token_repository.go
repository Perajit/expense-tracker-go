package auth

import "gorm.io/gorm"

type TokenRepository interface {
	WithTx(tx *gorm.DB) TokenRepository
	GetByTokenID(jti string) (*TokenEntity, error)
	Create(token *TokenEntity) error
	Revoke(token *TokenEntity) error
	RevokeByTokenID(jti string) error
	RevokeAllFromUser(userID uint) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenReposity(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) WithTx(tx *gorm.DB) TokenRepository {
	if tx == nil {
		return r
	}

	return &tokenRepository{db: tx}
}

func (r *tokenRepository) GetByTokenID(jti string) (*TokenEntity, error) {
	var token TokenEntity
	if err := r.db.Where("token_id = ?", jti).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepository) Create(token *TokenEntity) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) Revoke(token *TokenEntity) error {
	token.IsRevoked = true

	return r.db.Save(token).Error
}

func (r *tokenRepository) RevokeByTokenID(jti string) error {
	return r.db.Model(&TokenEntity{}).
		Where("token_id = ?", jti).
		Update("is_revoked", true).
		Error
}

func (r *tokenRepository) RevokeAllFromUser(userID uint) error {
	return r.db.Model(&TokenEntity{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).
		Error
}
