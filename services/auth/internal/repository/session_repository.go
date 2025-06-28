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

func (r *sessionRepository) GetSessionByID(ctx context.Context, sessionID uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error; err != nil {
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

func (r *sessionRepository) UpdateSession(ctx context.Context, updated *domain.Session) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", updated.ID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&session).Where("id = ?", updated.ID).Save(updated).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&session, "id = ?", updated.ID).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSessionByID(ctx context.Context, sessionID uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSessionByUserID(ctx context.Context, userID uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) RevokeSession(ctx context.Context, sessionID uint64) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&session).Where("id = ?", sessionID).Update("revoked", true).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}

	return &session, nil
}
