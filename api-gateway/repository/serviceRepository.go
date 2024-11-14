package repository

type ServiceRepository interface {
	GetDataByServiceName(serviceName string) (map[string]interface{}, bool)
}

type serviceRepository struct {
	dataService []map[string]interface{}
}

func NewServiceRepository() ServiceRepository {
	return &serviceRepository{
		dataService: []map[string]interface{}{
			{
				"serviceName": "user",
				"host":        "localhost:8000",
			},
		},
	}
}

func (sr *serviceRepository) GetDataByServiceName(serviceName string) (map[string]interface{}, bool) {
	for _, data := range sr.dataService {
		if data["serviceName"] == serviceName {
			return data, true
		}
	}
	return nil, false
}
