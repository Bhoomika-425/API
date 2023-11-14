package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique" validate:"required"`
	Location string `json:"location" validate:"required"`
	Field    string `json:"field" validate:"required"`
}

type Jobs struct {
	gorm.Model
	Company Company `json:"-" gorm:"ForeignKey:cid"`
	Cid     uint    `json:"cid" validate:"required"`
	// JiD              uint              `json:"jid" validate:"required"`
	Name             string            `json:"jobname" validate:"required"`
	Salary           string            `json:"salary" validate:"required"`
	NoticePeriod     string            `json:"notice_period" validate:"required"`
	MinNp            string            `json:"minNP" validate:"required"`
	MaxNP            string            `json:"maxNP" validate:"required"`
	Budget           string            `json:"budget" validate:"required"`
	Description      string            `json:"desc" validate:"required"`
	Minexp           string            `json:"minexp" validate:"required"`
	MaxMax           string            `json:"MaxMax" validate:"required"`
	Locations        []Location        `gorm:"many2many:job_locations;"`
	TechnologyStacks []TechnologyStack `gorm:"many2many:job_technology_stacks;"`
	WorkModes        []WorkMode        `gorm:"many2many:job_work_modes;"`
	Qualifications   []Qualification   `gorm:"many2many:job_qualifications;"`
	Shifts           []Shift           `gorm:"many2many:job_shifts;"`
	Jobtypes         []Jobtype         `gorm:"many2many:job_jobTypes;"`
}
type Location struct {
	gorm.Model
	Name string `json:"locname" gorm:"unique"`
}

type TechnologyStack struct {
	gorm.Model
	Skills string `json:"skills"  gorm:"unique"`
}
type WorkMode struct {
	gorm.Model
	Modename string `json:"modeName" gorm:"unique"`
}
type Qualification struct {
	gorm.Model
	Eligibility string `json:"eligibility"  gorm:"unique"`
}
type Shift struct {
	gorm.Model
	Shift string `json:"shift"  gorm:"unique"`
}
type Jobtype struct {
	gorm.Model
	Jobtype string `json:"jobType"  gorm:"unique"`
}

type NewJobRequest struct {
	Company      Company
	Jid          uint   `json:"jid" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Salary       string `json:"salary" validate:"required"`
	NoticePeriod string `json:"noticeperiod" validate:"required"`
	MinNp        string `json:"minnp" validate:"required"`
	MaxNP        string `json:"maxnp" validate:"required"`
	Budget       string `json:"budget" validate:"required"`
	Description  string `json:"desc" validate:"required"`
	Minexp       string `json:"minexp" validate:"required"`
	MaxMax       string `json:"maxexp" validate:"required"`
	Jobloc       []uint `json:"location" validate:"required"`
	Skills       []uint `json:"skills" validate:"required"`
	Mode         []uint `json:"mode" validate:"required"`
	Degree       []uint `json:"degree" validate:"required"`
	Shift        []uint `json:"shift" validate:"required"`
	Type         []uint `json:"type" validate:"required"`
}

type NewJobResponse struct {
	ID uint
}
type NewUserApplication struct {
	Name string       `json:"appname"`
	Age  string       `json:"age"`
	ID   uint         `json:"jid"`
	Jobs Requestfield `json:"Requestfield"`
}
type Requestfield struct {
	Name             string `json:"jobName" validate:"required"`
	NoticePeriod     uint   `json:"noticePeriod" validate:"required"`
	Locations        []uint `json:"location" `
	TechnologyStacks []uint `json:"technologyStack" `
	Experience       uint   `json:"experience" validate:"required"`
	Degree           []uint `json:"qualifications"`
	Shifts           []uint `json:"shifts"`
}

// type NewUserApplication struct {
// 	Name string       `json:"appname"`
// 	Age  string       `json:"age"`
// 	ID   uint         `json:"jid"`
// 	JobName         string `json:"jobName" validate:"required"`
// 	NoticePeriod string `json:"noticePeriod" validate:"required"`
// 	Jobloc       []uint `json:"location" `
// 	Skills       []uint `json:"technologyStack" `
// 	Experience   string `json:"experience" validate:"required"`
// 	Degree       []uint `json:"qualifications"`
// 	Shift        []uint `json:"shifts"`
// }
