package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ArchiveUserService struct {
	repo     repository.ArchiveUser
	Hub      *Hub
	Services *Services
}

func NewArchiveUserService(repo repository.ArchiveUser, hub *Hub) *ArchiveUserService {
	return &ArchiveUserService{repo: repo, Hub: hub}
}

func (s *ArchiveUserService) FindArchiveUser(filter *domain.ArchiveUserFilter) (domain.Response[domain.ArchiveUser], error) {
	return s.repo.FindArchiveUser(filter)
}

func (s *ArchiveUserService) CreateArchiveUser(ArchiveUserID string, ArchiveUser *domain.User) (*domain.ArchiveUser, error) {
	return s.repo.CreateArchiveUser(ArchiveUserID, ArchiveUser)
}

func (s *ArchiveUserService) DeleteArchiveUser(id string, ArchiveUserID string) (domain.ArchiveUser, error) {
	var result domain.ArchiveUser

	// delete images.
	allImages, err := s.Services.ArchiveImage.FindArchiveImage(&domain.ArchiveImageFilter{ServiceId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allImages.Data {
		// fmt.Println("Remove image: ", allImages.Data[i].ID)
		_, err = s.Services.ArchiveImage.DeleteArchiveImage(allImages.Data[i].ID.Hex())
	}

	// delete taskWorkers.
	allTaskWorkers, err := s.Services.ArchiveTaskWorker.FindArchiveTaskWorker(&domain.ArchiveTaskWorkerFilter{WorkerId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allTaskWorkers.Data {
		// fmt.Println("Remove TaskWorkers: ", allTaskWorkers.Data[i].ID)
		_, err = s.Services.ArchiveTaskWorker.DeleteArchiveTaskWorker(allTaskWorkers.Data[i].ID.Hex(), ArchiveUserID)
	}

	// delete workHistory.
	allWorkHistory, err := s.Services.ArchiveWorkHistory.FindArchiveWorkHistory(domain.ArchiveWorkHistoryFilter{WorkerId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allWorkHistory.Data {
		// fmt.Println("Remove WorkHistory: ", allWorkHistory.Data[i].ID)
		_, err = s.Services.ArchiveWorkHistory.DeleteArchiveWorkHistory(allWorkHistory.Data[i].ID.Hex(), ArchiveUserID)
	}

	// delete pay.
	allPay, err := s.Services.ArchivePay.FindArchivePay(&domain.ArchivePayFilter{WorkerId: []string{id}})
	if err != nil {
		return result, err
	}
	for i := range allPay.Data {
		// fmt.Println("Remove Pay: ", allPay.Data[i].ID)
		_, err = s.Services.Pay.DeletePay(allPay.Data[i].ID.Hex(), ArchiveUserID)
	}

	// delete notify.
	allNotify, err := s.Services.ArchiveNotify.FindArchiveNotifyPopulate(&domain.ArchiveNotifyFilter{UserTo: []*string{&id}})
	if err != nil {
		return result, err
	}
	for i := range allNotify.Data {
		// fmt.Println("Remove Notify: ", allNotify.Data[i].ID)
		_, err = s.Services.ArchiveNotify.DeleteArchiveNotify(allNotify.Data[i].ID.Hex(), ArchiveUserID)
	}

	result, err = s.repo.DeleteArchiveUser(id)

	return result, err
}
