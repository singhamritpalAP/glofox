package services

type IService interface {
	CreateClass(name, startDateStr, endDateStr string, capacity int) error
	BookClass(className, memberName, dateStr string) error
}
