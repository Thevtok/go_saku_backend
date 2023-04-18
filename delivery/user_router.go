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
	secretKey := []byte(utils.DotEnv("KEY"))

	db := config.LoadDatabase()
	defer db.Close()

	authMiddleware := controller.AuthMiddleware(secretKey)

	r := gin.Default()

	userRouter := r.Group("/user")
	userRouter.Use(authMiddleware)

	repo := repository.NewUserRepository(db)
	usecase := usecase.NewUserUseCase(repo)
	auth := controller.NewUserAuth(usecase)
	controller := controller.NewUserController(usecase)

	r.POST("/login", auth.Login)
	r.POST("/register", controller.Register)

	userRouter.GET("", controller.FindUsers)
	userRouter.GET("/:id", controller.FindUserByID)

	userRouter.PUT("", controller.Edit)
	userRouter.DELETE("/:id", controller.Unreg)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}

}
