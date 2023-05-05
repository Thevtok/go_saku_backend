package controller

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TxUseCaseMock struct {
	mock.Mock
}
type TxControllerTestSuite struct {
	suite.Suite
	routerMock   *gin.Engine
	bankCaseMock *BankAccUsecaseMock

	useCaseMock *UserUseCaseMock
}

func (suite *TxControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(UserUseCaseMock)
}
func TestTxController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
