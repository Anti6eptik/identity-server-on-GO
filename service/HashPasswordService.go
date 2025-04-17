package service


import "golang.org/x/crypto/bcrypt"


type HashPasswordService struct {}

func NewHashPasswordService() *HashPasswordService {
	return &HashPasswordService{}
}

func (srv HashPasswordService) CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (srv HashPasswordService) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}