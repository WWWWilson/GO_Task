package service

import (
	"golang_task3/models"
	"golang_task3/repository"
)

type StudentService struct {
    repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
    return &StudentService{repo: repo}
}

func (s *StudentService) CreateStudent(name string, age int, grade string) (*models.Student, error) {
    student := &models.Student{
        Name:  name,
        Age:   age,
        Grade: grade,
    }
    
    err := s.repo.CreateStudent(student)
    if err != nil {
        return nil, err
    }
    
    return student, nil
}

func (s *StudentService) GetStudent(id int) (*models.Student, error) {
    return s.repo.GetStudentByID(id)
}

func (s *StudentService) GetAllStudents() ([]*models.Student, error) {
    return s.repo.GetAllStudents()
}

func (s *StudentService) UpdateStudent( name string, age int, grade string) (int64,error) {
    student := &models.Student{
        Name:  name,
        Age:   age,
        Grade: grade,
    }
    
    return s.repo.UpdateStudent(student)
}

func (s *StudentService) DeleteStudent(age int) (int64, error) {
    return s.repo.DeleteStudent(age)
}

func (s *StudentService) GetStudentsByGrade(grade string) ([]*models.Student, error) {
    return s.repo.GetStudentsByGrade(grade)
}