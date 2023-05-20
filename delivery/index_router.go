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

	authMiddlewareRole := controller.AuthMiddlewareRole()

	r := gin.Default()

	// User Router
	userRouter := r.Group("/user")
	userRouter.Use(authMiddlewareUsername)

	// USER DEPEDENCY
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userAuth := controller.NewUserAuth(userUsecase)
	bankAccRepo := repository.NewBankAccRepository(db)
	bankAccusecase := usecase.NewBankAccUsecase(bankAccRepo)
	bankAccController := controller.NewBankAccController(bankAccusecase)
	cardRepo := repository.NewCardRepository(db)
	cardUsecase := usecase.NewCardUsecase(cardRepo)
	cardController := controller.NewCardController(cardUsecase)
	photoRepo := repository.NewPhotoRepository(db)
	photoUsecase := usecase.NewPhotoUseCase(photoRepo)
	photoController := controller.NewPhotoController(photoUsecase)
	userController := controller.NewUserController(userUsecase, bankAccusecase, cardUsecase, photoUsecase)
	authMiddlewareIdExist := userController.AuthMiddlewareIDExist()

	// USER GROUP
	r.POST("/login", userAuth.Login)
	r.POST("/register", userController.Register)

	r.GET("user", authMiddlewareRole, userController.FindUsers)
	userRouter.GET("username/:username", userController.FindUserByUsername)
	r.GET("user/:phone_number", userController.FindUserByPhone)
	// r.PUT("user/:user_id", controller.AuthMiddlewareID(), userController.Edit)
	r.PUT("user/pass/:user_id", authMiddlewareIdExist, userController.EditEmailPassword)
	r.PUT("user/profile/:user_id", authMiddlewareIdExist, userController.EditProfile)
	r.DELETE("user/:user_id", authMiddlewareIdExist, userController.Unreg)

	// Bank Accont Router
	bankAccRouter := r.Group("/user/bank")
	bankAccRouter.Use(authMiddlewareIdExist)

	// Bank Acc Depedency

	r.GET("user/bank", authMiddlewareRole, bankAccController.FindAllBankAcc)
	bankAccRouter.GET("/:user_id", bankAccController.FindBankAccByUserID)
	bankAccRouter.GET("/:user_id/:account_id", bankAccController.FindBankAccByAccountID)
	bankAccRouter.POST("/add/:user_id", bankAccController.CreateBankAccount)
	bankAccRouter.PUT("update/:user_id/:account_id", bankAccController.Edit)
	bankAccRouter.DELETE("/:user_id", bankAccController.UnregAll)
	bankAccRouter.DELETE("/:user_id/:account_id", bankAccController.UnregByAccountID)

	// CarDUnregByAccountID Router
	cardRouter := r.Group("/user/card")
	cardRouter.Use(authMiddlewareIdExist)

	// Card Depedency

	r.GET("user/card", authMiddlewareRole, cardController.FindAllCard)
	cardRouter.GET("/:user_id", cardController.FindCardByUserID)
	cardRouter.GET("/:user_id/:card_id", cardController.FindCardByCardID)
	cardRouter.POST("/add/:user_id", cardController.CreateCardID)
	cardRouter.PUT("/update/:user_id/:card_id", cardController.Edit)
	cardRouter.DELETE("/:user_id", cardController.UnregAll)
	cardRouter.DELETE("/:user_id/:card_id", cardController.UnregByCardID)

	// Photo Router
	photoRouter := r.Group("/user/photo")
	photoRouter.Use(authMiddlewareIdExist)

	// Photo Depedency

	photoRouter.POST("/:user_id", photoController.Upload)
	photoRouter.GET("/:user_id", photoController.Download)
	photoRouter.PUT("/:user_id", photoController.Edit)
	photoRouter.DELETE("/:user_id", photoController.Remove)

	//TX Router
	txRouter := r.Group("/user/tx")
	txRouter.Use(authMiddlewareIdExist)

	// TX Depedency
	txRepo := repository.NewTxRepository(db)
	txUsecase := usecase.NewTransactionUseCase(txRepo, userRepo)
	txController := controller.NewTransactionController(txUsecase, userUsecase, bankAccusecase, cardUsecase)

	txRouter.POST("/tf/:user_id", txController.CreateTransferTransaction)
	txRouter.POST("depo/bank/:user_id/:bank_account_id", txController.CreateDepositBank)
	txRouter.POST("depo/card/:user_id/:card_id", txController.CreateDepositCard)
	txRouter.POST("wd/:user_id/:bank_account_id", txController.CreateWithdrawal)
	txRouter.POST("redeem/:user_id/:pe_id", txController.CreateRedeemTransaction)
	txRouter.GET(":user_id", txController.GetTxBySenderId)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}
}
