package services

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockAPIClient struct {
	mock.Mock
}

const (
	// Mock CABA Transport - Parking Rules Response - Ok
	ParkingRulesResponseOK = `{
        "totalFull": 1,
        "instancias": [
            {
                "nombre": "Test Rule",
                "claseId": "1",
                "clase": "Test Class",
                "id": "123",
                "distancia": "100",
                "contenido": {
                    "contenido": [
                        {
                            "nombreId": "calle",
                            "nombre": "Calle",
                            "posicion": "1",
                            "valor": "Corrientes"
                        },
                        {
                            "nombreId": "altura",
                            "nombre": "Altura",
                            "posicion": "2",
                            "valor": "1000"
                        },
                        {
                            "nombreId": "permiso",
                            "nombre": "Permiso",
                            "posicion": "3",
                            "valor": "Permitido"
                        },
                        {
                            "nombreId": "horario",
                            "nombre": "Horario",
                            "posicion": "4",
                            "valor": "08:00-20:00"
                        },
                        {
                            "nombreId": "lado",
                            "nombre": "Lado",
                            "posicion": "5",
                            "valor": "Izquierdo"
                        }
                    ]
                }
            }
        ],
        "total": 1
    }`
)

func (m *MockAPIClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockAPIClient) ParkingRules(lat, long float64) (SimplifiedRules, error) {
	arg := m.Called(lat, long)
	if arg.Get(0) == nil {
		return nil, arg.Error(1)
	}
	return arg.Get(0).(SimplifiedRules), arg.Error(1)
}
