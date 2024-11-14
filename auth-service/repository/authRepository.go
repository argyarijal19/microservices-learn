package repository

type AuthRepository interface {
	GetDataByClientID(clientID string) (map[string]interface{}, bool)
}

type authRepository struct {
	dummyData []map[string]interface{}
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		dummyData: []map[string]interface{}{
			{
				"client_id":     "123",
				"client_secret": "456",
				"scope":         "read",
			},
			{
				"client_id":     "789",
				"client_secret": "abc",
				"scope":         "write",
			},
		},
	}
}

func (ar *authRepository) GetDataByClientID(clientID string) (map[string]interface{}, bool) {
	for _, data := range ar.dummyData {
		if data["client_id"] == clientID {
			return data, true
		}
	}
	return nil, false
}
