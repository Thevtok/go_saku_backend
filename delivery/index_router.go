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

	authMiddlewareUsername := controller.AuthMiddleware()
	authMiddlewareId := controller.AuthMiddlewareID()

	r := gin.Default()

	// User Router
	userRouter := r.Group("/user")
	userRouter.Use(authMiddlewareUsername)
	// USER DEPEDENCY

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userAuth := controller.NewUserAuth(userUsecase)
	userController := controller.NewUserController(userUsecase)

	r.POST("/login", userAuth.Login)
	r.POST("/register", userController.Register)
	// USER GROUP

	userRouter.GET("", userController.FindUsers)
	userRouter.GET("/:username", userController.FindUserByUsername)
	r.PUT("user/:user_id", controller.AuthMiddlewareID(), userController.Edit)
	userRouter.DELETE("/:username", userController.Unreg)

	// Bank Accont Router
	bankAccRouter := r.Group("/user/bank")
	bankAccRouter.Use(authMiddlewareId)

	bankAccRepo := repository.NewBankAccRepository(db)
	bankAccusecase := usecase.NewBankAccUsecase(bankAccRepo)
	bankAccController := controller.NewBankAccController(bankAccusecase)

	bankAccRouter.GET("", bankAccController.FindAllBankAcc)
	bankAccRouter.GET("/:user_id", bankAccController.FindBankAccByID)
	bankAccRouter.POST("/add/:user_id", bankAccController.Register)
	bankAccRouter.PUT("update/:user_id", bankAccController.Edit)

	bankAccRouter.DELETE("/:user_id", bankAccController.Unreg)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}

}
