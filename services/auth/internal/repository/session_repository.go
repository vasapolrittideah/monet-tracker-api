package repository

import (
	"context"

	auth "github.com/vasapolrittideah/money-tracker-api/services/auth/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) auth.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "token = ? AND revoked = false", token).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) UpdateSession(
	ctx context.Context,
	id uint64,
	updates map[string]any,
) (*domain.Session, error) {
	if err := r.db.WithContext(ctx).Model(&domain.Session{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSessionByID(ctx context.Context, id uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSessionByUserID(ctx context.Context, userID uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) RevokeSession(ctx context.Context, id uint64) (*domain.Session, error) {
	if err := r.db.WithContext(ctx).Model(&domain.Session{}).Where("id = ?", id).Update("revoked", true).Error; err != nil {
		return nil, err
	}

	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &session, nil
}
