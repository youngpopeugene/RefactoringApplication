package models

import (
	"time"
)

type Substation struct {
	Name               string `gorm:"primaryKey" json:"name"`
	Location           string `json:"location"`
	YearOfConstruction int    `json:"year_of_construction"`
	CommissioningYear  int    `json:"commissioning_year"`
}

type Factory struct {
	Name string `gorm:"primaryKey" json:"name"`
	City string `json:"city"`
}

type RangeOfHighVoltageEquipment struct {
	ID                  int    `gorm:"primaryKey" json:"id"`
	HighVoltageSwitch   string `json:"high_voltage_switch"`
	MediumVoltageSwitch string `json:"medium_voltage_switch"`
	LowVoltageSwitch    string `json:"low_voltage_switch"`
}

type CableLine struct {
	Mark string `gorm:"primaryKey" json:"mark"`
}

type TireSection struct {
	Name string `gorm:"primaryKey" json:"name"`
}

type CellKVL struct {
	DispatchName                string `gorm:"primaryKey" json:"dispatch_name"`
	CableLines                  string `gorm:"references:CableLine" json:"cable_lines"`
	CurrentTransformer          string `json:"current_transformer"`
	Switch                      string `json:"switch"`
	ProtectionTransformer       string `json:"protection_transformer"`
	TireSection                 string `gorm:"references:TireSection" json:"tire_section"`
	NumberOfCurrentTransformers int    `json:"number_of_current_transformers"`
}

type Fuse struct {
	Mark string `gorm:"primaryKey" json:"mark"`
}

type CellTN struct {
	DispatchName       string `gorm:"primaryKey" json:"dispatch_name"`
	VoltageTransformer string `json:"voltage_transformer"`
	Fuse               string `gorm:"references:Fuse" json:"fuse"`
	TireSection        string `gorm:"references:TireSection" json:"tire_section"`
}

type CellTSN struct {
	DispatchName         string `gorm:"primaryKey" json:"dispatch_name"`
	AuxiliaryTransformer string `json:"auxiliary_transformer"`
	Fuse                 string `gorm:"references:Fuse" json:"fuse"`
	TireSection          string `gorm:"references:TireSection" json:"tire_section"`
}

type NSS struct {
	ID             int `gorm:"primaryKey" json:"id"`
	RatedVoltageKV int `json:"rated_voltage_kv"`
}

type RangeOfStandardVoltage struct {
	ID                      int `gorm:"primaryKey" json:"id"`
	RatedWindingVoltageHVKV int `json:"rated_winding_voltage_hv_kv"`
	RatedWindingVoltageMVKV int `json:"rated_winding_voltage_mv_kv"`
	RatedWindingVoltageLVKV int `json:"rated_winding_voltage_lv_kv"`
}

type TypeOfTransformer struct {
	Type                   string `gorm:"primaryKey" json:"type"`
	PowerMVA               int    `json:"power_mva"`
	CoolingSystemType      string `json:"cooling_system_type"`
	RangeOfStandardVoltage int    `gorm:"references:RangeOfStandardVoltage" json:"range_of_standard_voltage"`
}

type Transformer struct {
	FactoryNumber               int       `gorm:"primaryKey" json:"factory_number"`
	NSS                         int       `gorm:"references:NSS" json:"nss"`
	Substation                  string    `gorm:"references:Substation" json:"substation"`
	Factory                     string    `gorm:"references:Factory" json:"factory"`
	Type                        string    `gorm:"references:TypeOfTransformer" json:"type"`
	DateOfManufacture           time.Time `json:"date_of_manufacture"`
	CommissioningDate           time.Time `json:"commissioning_date"`
	DispatchName                string    `json:"dispatch_name"`
	RangeOfHighVoltageEquipment int       `gorm:"references:RangeOfHighVoltageEquipment" json:"range_of_high_voltage_equipment"`
	TireSection                 string    `json:"tire_section"`
}

type Request struct {
	ID                       int       `gorm:"primaryKey" json:"id"`
	TransformerFactoryNumber int       `gorm:"references:Transformer" json:"transformer_factory_number"`
	WorkerUsername           string    `gorm:"references:User" json:"worker_username"`
	IsCompleted              bool      `json:"is_completed"`
	DateOpened               time.Time `json:"date_opened"`
	DateClosed               time.Time `json:"date_closed"`
}

type Role string

const (
	RoleWorker     Role = "WORKER"
	RoleDispatcher Role = "DISPATCHER"
)

type User struct {
	Username string `gorm:"primaryKey" json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}
