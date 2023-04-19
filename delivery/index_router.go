package delivery

import (
	"log"

	"github.com/ReygaFitra/inc-final-project.git/config"
	"github.com/ReygaFitra/inc-final-project.git/controller"
	"github.com/ReygaFitra/inc-final-project.git/repository"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func RunServer() {

	db := config.LoadDatabase()
	defer db.Close()

	authMiddleware := controller.AuthMiddleware()

	r := gin.Default()

	// User Router
	userRouter := r.Group("/user")
	userRouter.Use(authMiddleware)
	// USER DEPEDENCY

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userAuth := controller.NewUserAuth(userUsecase)
	userController := controller.NewUserController(userUsecase)

	r.POST("/login", userAuth.Login)
	r.POST("/register", userController.Register)
	// USER GROUP

	userRouter.GET("", userController.FindUsers)
	userRouter.GET("/:id", userController.FindUserByID)
	userRouter.PUT("", userController.Edit)
	userRouter.DELETE("/:id", userController.Unreg)

	// Bank Accont Router
	bankAccRouter := r.Group("/user/bank")
	bankAccRouter.Use(authMiddleware)

	bankAccRepo := repository.NewBankAccRepository(db)
	bankAccusecase := usecase.NewBankAccUsecase(bankAccRepo)
	bankAccController := controller.NewBankAccController(bankAccusecase)

	bankAccRouter.GET("", bankAccController.FindAllBankAcc)
	bankAccRouter.GET("/:userID/:accountID", bankAccController.FindBankAccByID)
	bankAccRouter.POST("/add", bankAccController.Register)
	bankAccRouter.PUT("update", bankAccController.Edit)

	bankAccRouter.DELETE("/:id", bankAccController.Unreg)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}

}
