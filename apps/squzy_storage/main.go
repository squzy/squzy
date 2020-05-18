package squzy_storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"squzy/apps/squzy_storage/application"
	"squzy/apps/squzy_storage/config"
	"squzy/apps/squzy_storage/server"
	"squzy/internal/database"
)

func main()  {
	cnfg := config.New()
	db, err := database.New(func() (db *gorm.DB, e error) {
		return gorm.Open(
			"postgres",
			fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s connect_timeout=10 sslmode=disable",
				cnfg.GetDbHost(),
				cnfg.GetDbPort(),
				cnfg.GetDbName(),
				cnfg.GetDbUser(),
				cnfg.GetDbPassword(),
			))
	})
	if err != nil {
		log.Fatal(err)
	}
	apiService := application.NewService(db)
	storageServ := server.NewServer(cnfg, apiService)
	log.Fatal(storageServ.Run())
}