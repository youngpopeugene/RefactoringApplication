package gateway

import (
	"app/internal/domain"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

//Get by primary key

func (gw *gateway) GetSubstationByPK(pk string) (domain.Substation, error) {
	var obj domain.Substation
	if tx := gw.db.First(&obj, "name = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("substation not found")
	}
	return obj, nil
}

func (gw *gateway) GetFactoryByPK(pk string) (domain.Factory, error) {
	var obj domain.Factory
	if tx := gw.db.First(&obj, "name = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("factory not found")
	}
	return obj, nil
}

func (gw *gateway) GetRangeOfHighVoltageEquipmentByPK(pk string) (domain.RangeOfHighVoltageEquipment, error) {
	pkInt, _ := strconv.Atoi(pk)
	var obj domain.RangeOfHighVoltageEquipment
	if tx := gw.db.First(&obj, "id = ?", pkInt); tx.RowsAffected == 0 {
		return obj, errors.New("range of high voltage equipment not found")
	}
	return obj, nil
}

func (gw *gateway) GetCableLineByPK(pk string) (domain.CableLine, error) {
	var obj domain.CableLine
	if tx := gw.db.First(&obj, "mark = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("cable line not found")
	}
	return obj, nil
}

func (gw *gateway) GetTireSectionByPK(pk string) (domain.TireSection, error) {
	var obj domain.TireSection
	if tx := gw.db.First(&obj, "name = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("tire section not found")
	}
	return obj, nil
}

func (gw *gateway) GetCellKVLByPK(pk string) (domain.CellKVL, error) {
	var obj domain.CellKVL
	if tx := gw.db.First(&obj, "dispatch_name = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("cell kvl not found")
	}
	return obj, nil
}

func (gw *gateway) GetFuseByPK(pk string) (domain.Fuse, error) {
	var obj domain.Fuse
	if tx := gw.db.First(&obj, "mark = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("fuse not found")
	}
	return obj, nil
}

func (gw *gateway) GetCellTNByPK(pk string) (domain.CellTN, error) {
	var obj domain.CellTN
	if tx := gw.db.First(&obj, "dispatch_name = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("cell tnb not found")
	}
	return obj, nil
}

func (gw *gateway) GetCellTSNByPK(pk string) (domain.CellTSN, error) {
	var obj domain.CellTSN
	if tx := gw.db.First(&obj, "dispatch_name = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("cell tsn not found")
	}
	return obj, nil
}

func (gw *gateway) GetNSSByPK(pk string) (domain.NSS, error) {
	pkInt, _ := strconv.Atoi(pk)
	var obj domain.NSS
	if tx := gw.db.First(&obj, "id = ?", pkInt); tx.RowsAffected == 0 {
		return obj, errors.New("nss not found")
	}
	return obj, nil
}

func (gw *gateway) GetRangeOfStandardVoltageByPK(pk string) (domain.RangeOfStandardVoltage, error) {
	pkInt, _ := strconv.Atoi(pk)
	var obj domain.RangeOfStandardVoltage
	if tx := gw.db.First(&obj, "id = ?", pkInt); tx.RowsAffected == 0 {
		return obj, errors.New("range of standard voltage not found")
	}
	return obj, nil
}

func (gw *gateway) GetTypeOfTransformerByPK(pk string) (domain.TypeOfTransformer, error) {
	var obj domain.TypeOfTransformer
	if tx := gw.db.First(&obj, "type = ?", pk); tx.RowsAffected == 0 {
		return obj, errors.New("type of transformer not found")
	}
	return obj, nil
}

func (gw *gateway) GetTransformerByPK(pk string) (domain.Transformer, error) {
	pkInt, _ := strconv.Atoi(pk)
	var obj domain.Transformer
	if tx := gw.db.First(&obj, "factory_number = ?", pkInt); tx.RowsAffected == 0 {
		return obj, errors.New("transformer not found")
	}
	return obj, nil
}

func (gw *gateway) GetRequestByPK(pk string) (domain.Request, error) {
	pkInt, _ := strconv.Atoi(pk)
	var obj domain.Request
	if tx := gw.db.First(&obj, "id = ?", pkInt); tx.RowsAffected == 0 {
		return obj, errors.New("request not found")
	}
	return obj, nil
}

//Get all users workers

func (gw *gateway) GetAllUsersWorkers() []domain.User {
	role := domain.RoleWorker
	var users []domain.User
	gw.db.Where("role = ?", role).Find(&users)
	return users
}

//Get all unique location

func (gw *gateway) GetAllLocations() []string {
	var locations []string
	gw.db.Model(&domain.Substation{}).Select("DISTINCT location").Find(&locations)
	return locations
}

//Get transformers by location of substation

func (gw *gateway) GetTransformersByLocation(location string) []domain.Transformer {
	var transformers []domain.Transformer
	gw.db.Raw("SELECT * "+
		"FROM transformers "+
		"JOIN substations on transformers.substation = substations.name "+
		"WHERE substations.location = @location",
		sql.Named("location", location)).Find(&transformers)
	return transformers
}

//Get all requests

func (gw *gateway) GetAllRequests() []domain.Request {
	var requests []domain.Request
	gw.db.Find(&requests)
	return requests
}

//Get requests by worker_username

func (gw *gateway) GetRequestsByWorkerUsername(workerUsername string) []domain.Request {
	var requests []domain.Request
	gw.db.Where("worker_username = ?", workerUsername).Find(&requests)
	return requests
}

//Create request

func (gw *gateway) CreateRequest(workerUsername string, transformerFactoryNumber int) (domain.Request, error) {
	var transformer domain.Transformer
	gw.db.First(&transformer, transformerFactoryNumber)

	var user domain.User
	gw.db.First(&user, workerUsername)

	if user.Role != domain.RoleWorker {
		return domain.Request{}, errors.New("only worker's username could be passed")
	}

	request := domain.Request{
		TransformerFactoryNumber: transformer.FactoryNumber,
		WorkerUsername:           user.Username,
		IsCompleted:              false,
		DateOpened:               time.Now(),
	}

	gw.db.Model(&transformer).Association("TransformerFactoryNumber").Append(&request)
	gw.db.Model(&user).Association("WorkerUsername").Append(&request)
	gw.db.Create(&request)
	return request, nil
}

//Update request (close it)

func (gw *gateway) UpdateRequest(workerUsername string, transformerFactoryNumber int) (domain.Request, error) {

	var transformer domain.Transformer
	gw.db.First(&transformer, transformerFactoryNumber)

	var user domain.User
	gw.db.First(&user, workerUsername)

	if user.Role != domain.RoleWorker {
		return domain.Request{}, errors.New("only worker's username could be passed")
	}

	request := domain.Request{
		TransformerFactoryNumber: transformer.FactoryNumber,
		WorkerUsername:           user.Username,
	}
	gw.db.First(&request)
	request.IsCompleted = true
	request.DateClosed = time.Now()

	gw.db.Model(&transformer).Association("TransformerFactoryNumber").Append(&request)
	gw.db.Model(&user).Association("WorkerUsername").Append(&request)
	gw.db.Save(&request)
	return request, nil
}

func (gw *gateway) SaveUser(user domain.User) error {
	if err := gw.db.Create(&user).Error; err != nil {
		return errors.New("user not saved")
	}
	return nil
}

func (gw *gateway) GetUserByName(username string) (domain.User, error) {
	var user domain.User
	if tx := gw.db.Where("username = ?", username).First(&user); tx.RowsAffected == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (gw *gateway) SaveRefreshToken(refreshToken domain.Token) error {
	if err := gw.db.Create(&refreshToken).Error; err != nil {
		return err
	}
	return nil
}

func (gw *gateway) DeactivateRefreshTokens(username string) error {
	if err :=
		gw.db.Model(&domain.Token{}).
			Where("username = ?", username).
			Update("active", false).Error; err != nil {
		return err
	}
	return nil
}
