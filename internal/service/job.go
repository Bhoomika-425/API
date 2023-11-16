package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"project/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func (s *Service) ViewJobById(ctx context.Context, jid uint64) (models.Jobs, error) {
	jobData, err := s.UserRepo.Jobbyjid(ctx, jid)
	if err != nil {
		return models.Jobs{}, nil
	}
	return jobData, nil
}

func (s *Service) ViewAllJobs(ctx context.Context) ([]models.Jobs, error) {
	jobDatas, err := s.UserRepo.FetchAllJobs(ctx)
	if err != nil {
		return nil, err
	}
	return jobDatas, nil

}
func (s *Service) ViewJob(ctx context.Context, cid uint64) ([]models.Jobs, error) {
	jobData, err := s.UserRepo.Jobbycid(ctx, cid)
	if err != nil {
		return nil, err
	}
	return jobData, nil
}

func (s *Service) AddJobDetails(ctx context.Context, jobRequest models.NewJobRequest, cid uint) (models.NewJobResponse, error) {
	// jobData.Cid = uint(cid)
	requestJob := models.Jobs{
		Cid:          cid,
		Name:         jobRequest.Name,
		Salary:       jobRequest.Salary,
		NoticePeriod: jobRequest.NoticePeriod,
		MinNp:        jobRequest.MinNp,
		MaxNP:        jobRequest.MaxNP,
		Budget:       jobRequest.Budget,
		Description:  jobRequest.Description,
		Minexp:       jobRequest.Minexp,
		MaxMax:       jobRequest.MaxMax,
	}
	for _, v := range jobRequest.Jobloc {
		tempData := models.Location{
			Model: gorm.Model{
				ID: v,
			},
		}
		requestJob.Locations = append(requestJob.Locations, tempData)
	}
	for _, v := range jobRequest.Skills {
		tempData := models.TechnologyStack{
			Model: gorm.Model{
				ID: v,
			},
		}
		requestJob.TechnologyStacks = append(requestJob.TechnologyStacks, tempData)
	}
	for _, v := range jobRequest.Mode {
		tempData := models.WorkMode{
			Model: gorm.Model{
				ID: v,
			},
		}
		requestJob.WorkModes = append(requestJob.WorkModes, tempData)
	}

	for _, v := range jobRequest.Degree {
		tempData := models.Qualification{
			Model: gorm.Model{
				ID: v,
			},
		}
		requestJob.Qualifications = append(requestJob.Qualifications, tempData)
	}
	for _, v := range jobRequest.Shift {
		tempData := models.Shift{
			Model: gorm.Model{
				ID: v,
			},
		}
		requestJob.Shifts = append(requestJob.Shifts, tempData)
	}
	for _, v := range jobRequest.Type {
		tempData := models.Jobtype{
			Model: gorm.Model{
				ID: v,
			},
		}
		requestJob.Jobtypes = append(requestJob.Jobtypes, tempData)
	}
	job, err := s.UserRepo.CreateUserJob(ctx, requestJob)
	if err != nil {
		return models.NewJobResponse{}, err
	}
	return job, nil
}

func (s *Service) ApplyJobs(ctx context.Context, applications []models.NewUserApplication) ([]models.NewUserApplication, error) {
	var wg = new(sync.WaitGroup)
	ch := make(chan models.NewUserApplication)
	var finalData []models.NewUserApplication
	for _, v := range applications {
		wg.Add(1)
		go func(v models.NewUserApplication) {
			defer wg.Done()
			var jobdata models.Jobs
			val, err := s.rdb.GettingCacheData(ctx, v.ID)
			if err == redis.Nil {
				dbData, err := s.UserRepo.CreateApplication(ctx, v.ID)
				if err != nil {
					return
				}
				err = s.rdb.AddingtoCache(ctx, v.ID, dbData)
				if err != nil {
					return
				}
				jobdata = dbData
			} else {
				err = json.Unmarshal([]byte(val), &jobdata)
				if err == redis.Nil {
					return
				}
				if err != nil {
					return
				}
			}
			check, v, err := s.CompareAndCheck(ctx, v, jobdata)
			if err != nil {
				return
			}
			if check {
				ch <- v
			}
			//finalData = append(finalData, v)
		}(v)
	}

	// check, v, err := s.CompareAndCheck(ctx, v)

	// if err != nil {
	// 	return nil, err
	// }
	// if check {
	// 	finalData = append(finalData, v)
	// }
	go func() {
		wg.Wait()
		close(ch)
	}()
	for v := range ch {
		finalData = append(finalData, v)
	}
	return finalData, nil
}

func (s *Service) CompareAndCheck(ctx context.Context, appData models.NewUserApplication, val models.Jobs) (bool, models.NewUserApplication, error) {
	// jobData, err := s.UserRepo.CreateApplication(ctx, appData.ID)
	// if err != nil {
	// 	return false, models.NewUserApplication{}, err
	// }

	// Define a variable to keep track of the number of matching conditions
	count := 0

	// Define the total number of conditions to match
	totalConditions := 7 // Adjust this number as per your specific conditions

	a, err := strconv.Atoi(val.Minexp)
	if err != nil {
		panic("parsing error")
	}
	if appData.Jobs.Experience >= uint(a) {
		count++
	}

	// if appData.Jobs.Salary == jobData.Salary {
	// 	count++
	// }
	b, err := strconv.Atoi(val.MinNp)
	if err != nil {
		panic("parsing error")
	}

	if appData.Jobs.NoticePeriod >= uint(b) {
		count++
	}

	// Compare job locations
	for _, v := range appData.Jobs.Locations {
		for _, v1 := range val.Locations {
			if v == v1.ID {
				count++
				break
			}
		}
	}

	// Compare qualifications
	for _, v := range appData.Jobs.Degree {
		for _, v1 := range val.Qualifications {
			if v == v1.ID {
				count++
				break
			}
		}
	}

	// Compare skills
	for _, v := range appData.Jobs.TechnologyStacks {
		for _, v1 := range val.TechnologyStacks {
			if v == v1.ID {
				count++
				break
			}
		}
	}

	// Compare shifts
	for _, v := range appData.Jobs.Shifts {
		for _, v1 := range val.Shifts {
			if v == v1.ID {
				count++
				break
			}
		}
	}

	// Check if at least 50% of the conditions match
	if count*2 >= totalConditions {
		return true, appData, nil
	}

	// If less than 50% match, return an empty struct
	return false, models.NewUserApplication{}, nil
}

//func (s *Service) compareAndCheck(ctx context.Context, applicationData models.NewUserApplication) (bool, models.NewUserApplication, error) {

// jobData, err := s.UserRepo.CreateApplication(ctx, applicationData.ID)
// if err != nil {
// 	return false, models.NewUserApplication{}, err
// }

// val := jobData

// if applicationData.Experience < val.Minexp {
// 	return false, models.NewUserApplication{}, nil
// }
// if applicationData.NoticePeriod < val.NoticePeriod {
// 	return false, models.NewUserApplication{}, nil
// }
// count := 0
// for _, v := range applicationData.Locations {
// 	for _, v1 := range val.Locations {
// 		if v == v1.ID {
// 			count++
// 		}
// 	}
// }
// if count == 0 {
// 	return false, models.NewUserApplication{}, nil
// }
// count = 0
// for _, v := range applicationData.Degree {
// 	for _, v1 := range val.Qualifications {
// 		if v == v1.ID {
// 			count++
// 		}
// 	}
// }
// if count == 0 {
// 	return false, models.NewUserApplication{}, nil
// }
// count = 0
// for _, v := range applicationData.TechnologyStacks {
// 	for _, v1 := range val.TechnologyStacks {
// 		if v == v1.ID {
// 			count++
// 		}
// 	}
// }
// if count == 0 {
// 	return false, models.NewUserApplication{}, nil
// }
// count = 0
// for _, v := range applicationData.Shifts {
// 	for _, v1 := range val.Shifts {
// 		if v == v1.ID {
// 			count++
// 		}
// 	}
// }
// if count == 0 {
// 	return false, models.NewUserApplication{}, nil
// }

// return true, applicationData, nil
// }
