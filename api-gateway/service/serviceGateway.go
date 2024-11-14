package service

import "api-gateway/repository"

type ServiceGateway interface {
	GetDataByServiceName(serviceName string) (map[string]interface{}, bool)
}

type serviceGateway struct {
	repos repository.ServiceRepository
}

func NewServiceGateway(r repository.ServiceRepository) ServiceGateway {
	return &serviceGateway{
		repos: r,
	}
}

func (s *serviceGateway) GetDataByServiceName(serviceName string) (map[string]interface{}, bool) {
	data, ok := s.repos.GetDataByServiceName(serviceName)
	if data == nil || !ok {
		return nil, false
	}
	return data, true
}
