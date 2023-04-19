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

	userRouter := r.Group("/user")
	userRouter.Use(authMiddleware)
	// USER DEPEDENCY

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userAuth := controller.NewUserAuth(userUsecase)
	userController := controller.NewUserController(userUsecase)

	// BANK DEPEDENCY

	bankRepo := repository.NewBankAccRepository(db)
	bankUsecase := usecase.NewBankAccUsecase(bankRepo)
	bankController := controller.NewBankAccController(bankUsecase)

	r.POST("/login", userAuth.Login)
	r.POST("/register", userController.Register)
	// USER GROUP
	userRouter.GET("", userController.FindUsers)
	userRouter.GET("/:id", userController.FindUserByID)

	userRouter.PUT("", userController.Edit)
	userRouter.DELETE("/:id", userController.Unreg)

	// BANK GROUP
	bankRouter := r.Group("/user/bank")
	bankRouter.Use(authMiddleware)

	bankRouter.GET("", bankController.FindAllBankAcc)
	bankRouter.GET("/:id", bankController.FindBankAccByID)
	bankRouter.POST("/add", bankController.Register)
	bankRouter.PUT("/update", bankController.Edit)
	bankRouter.DELETE("/delete/:id", bankController.Unreg)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}

}
