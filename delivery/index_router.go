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
	authMiddlewareRole := controller.AuthMiddlewareRole()

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

	r.GET("user", authMiddlewareRole, userController.FindUsers)
	userRouter.GET("/:username", userController.FindUserByUsername)
	// r.PUT("user/:user_id", controller.AuthMiddlewareID(), userController.Edit)
	r.PUT("user/:user_id", authMiddlewareId, userController.Edit)
	userRouter.DELETE("/:username", userController.Unreg)

	// Bank Accont Router
	bankAccRouter := r.Group("/user/bank")
	bankAccRouter.Use(authMiddlewareId)

	bankAccRepo := repository.NewBankAccRepository(db)
	bankAccusecase := usecase.NewBankAccUsecase(bankAccRepo)
	bankAccController := controller.NewBankAccController(bankAccusecase)

	r.GET("user/bank", authMiddlewareRole, bankAccController.FindAllBankAcc)
	bankAccRouter.GET("/:user_id", bankAccController.FindBankAccByUseerID)
	bankAccRouter.POST("/add/:user_id", bankAccController.CreateBankAccount)
	bankAccRouter.PUT("update/:user_id/:account_id", bankAccController.Edit)
	bankAccRouter.DELETE("/:user_id", bankAccController.UnregAll)
	bankAccRouter.DELETE("/:user_id/:account_id", bankAccController.UnregByAccountId)

	// Card Router
	cardRouter := r.Group("/user/card")
	cardRouter.Use(authMiddlewareId)

	cardRepo := repository.NewCardRepository(db)
	cardUsecase := usecase.NewCardUsecase(cardRepo)
	cardController := controller.NewCardController(cardUsecase)

	cardRouter.GET("user/card", cardController.FindAllCard)
	cardRouter.GET("/:user_id", cardController.FindCardByUserID)
	cardRouter.POST("/add/:user_id", cardController.CreateCardID)
	cardRouter.PUT("/update/:user_id/:card_id", cardController.Edit)
	cardRouter.DELETE("/:user_id", cardController.UnregAll)
	cardRouter.DELETE("/:user_id/card_id", cardController.UnregByCardId)

	// Photo Router
	photoRouter := r.Group("/user/photo")
	photoRouter.Use(authMiddlewareId)

	photoRepo := repository.NewPhotoRepository(db)
	photoUsecase := usecase.NewPhotoUseCase(photoRepo)
	photoController := controller.NewPhotoController(photoUsecase)

	photoRouter.POST("/:user_id", photoController.Upload)
	photoRouter.GET("/:user_id", photoController.Download)
	photoRouter.PUT("/:user_id", photoController.Edit)
	photoRouter.DELETE("/:user_id", photoController.Remove)

	//TX Router
	txRouter := r.Group("/user/tx")
	txRouter.Use(authMiddlewareId)

	txRepo := repository.NewTxRepository(db)
	txUsecase := usecase.NewTransactionUseCase(txRepo, userRepo)
	txController := controller.NewTransactionController(txUsecase, userUsecase)

	txRouter.POST("/tf/:user_id", txController.CreateTransferTransaction)
	txRouter.POST("depo/bank/:user_id", txController.CreateDepositBank)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}

}
