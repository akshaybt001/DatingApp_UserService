package helper

import (
	"fmt"
	"time"

	helperstruct "github.com/akshaybt001/DatingApp_UserService/entities/helperStruct"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashedPassword(hashedPass, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func CalculateAge(dob time.Time) int {
	currentDate := time.Now()

	years := currentDate.Year() - dob.Year()

	if currentDate.Month() < dob.Month() || currentDate.Month() == dob.Month() && currentDate.Day() < dob.Day() {
		years--
	}
	return years
}
func SearchForInterest(a []string, num string, beg int, end int) bool {
	for beg <= end {
		mid := (end + beg) / 2
		if a[mid] == num {
			return true
		} else if a[mid] < num {
			beg = mid + 1
		} else {
			end = mid - 1
		}
	}
	return false
}

func Abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}


func QuickSort(a []float64, b []helperstruct.Home, start, end int) {
	if start < end {
		p := partition(a, b, start, end)
		QuickSort(a, b, start, p-1)
		QuickSort(a, b, p+1, end)
	}
}
func partition(a []float64, b []helperstruct.Home, start, end int) int {
	pivot := a[end]
	i := start - 1
	for j := start; j < end; j++ {
		if a[j] < pivot {
			i++
			a[i], a[j] = a[j], a[i]
			b[i], b[j] = b[j], b[i]
		}
	}
	a[i+1], a[end] = a[end], a[i+1]
	b[i+1], b[end] = b[end], b[i+1]
	return i + 1
}
