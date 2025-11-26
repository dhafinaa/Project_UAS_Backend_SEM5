package service

import (
	"errors"
	"PROJECT_UAS/app/repository"
	"PROJECT_UAS/helper"
)

type AuthService struct {
	UserRepo       *repository.UserRepository
	PermissionRepo *repository.PermissionRepository
}

type LoginResponse struct {
	Token       string      `json:"token"`
	Refresh     string      `json:"refresh"`
	User        interface{} `json:"user"`
	Permissions []string    `json:"permissions"`
}

func (s *AuthService) Login(login string, password string) (*LoginResponse, error) {
	// 1. Ambil user berdasarkan username/email
	user, err := s.UserRepo.FindByLogin(login)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// 2. Cek status aktif
	if !user.Active_status {
		return nil, errors.New("akun tidak aktif")
	}

	// 3. Validasi password
	if !helper.CheckPasswordHash(password, user.Password_hash) {
		return nil, errors.New("password salah")
	}

	// 4. Ambil permissions berdasarkan role
	perms, err := s.PermissionRepo.GetByRole(user.Role_id)
	if err != nil {
		perms = []string{}
	}

	// 5. Generate JWT Access & Refresh Token
	token, refresh := helper.GenerateTokens(
		user.ID,
		user.Role_id,
		perms,
	)

	return &LoginResponse{
		Token:       token,
		Refresh:     refresh,
		User:        user,
		Permissions: perms,
	}, nil
}
