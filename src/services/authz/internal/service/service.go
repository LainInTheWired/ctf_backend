package service

type authzService struct {
}

type AuthzService interface {
}

func NewUserService(userrepository repository.UserRepository, redisrepository repository.RedisRepository) UserService {
	return &userService{
		usrepo: userrepository,
		rerepo: redisrepository,
	}
}
