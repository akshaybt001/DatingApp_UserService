package adapters

import (
	"github.com/akshaybt001/DatingApp_UserService/entities"
	helperstruct "github.com/akshaybt001/DatingApp_UserService/entities/helperStruct"
)

type AdapterInterface interface {
	UserSignup(entities.User) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetUserByPhone(phone string) (entities.User, error)
	GetAdminByEmail(email string) (entities.Admin, error)
	CreateProfile(userID string) error
	GetProfileIdByUserId(userId string) (string, error)

	AdminAddInterest(entities.Interests) error
	GetInterestByName(interest string) (entities.Interests, error)
	AdminUpdateInterest(entities.Interests) error
	AdminGetAllInterest() ([]entities.Interests, error)
	GetGenderByName(gender string) (entities.Gender, error)
	AdminAddGender(entities.Gender) error
	AdminUpdateGender(entities.Gender) error
	AdminGetAllGender() ([]entities.Gender, error)

	UserAddInterest(interest entities.UserInterests) error
	UserDeleteInterest(interest entities.UserInterests) error
	UserGetAllInterest(profileId string) ([]entities.Interests, error)
	GetInterestById(id int) (helperstruct.InterestHelper, error)
	GetUserInterestById(profileId string, interestId int) (entities.UserInterests, error)
	UserAddAddress(entities.Address) error
	UserEditAddress(entities.Address) error
	GetAddressByProfileId(profileId string) (entities.Address, error)
	GetGenderById(id int) (helperstruct.GenderHelper, error)
	GetGenderByProfileId(id string) (entities.UserGenders, error)
	GetUserGenderById(profileId string, genderId int) (entities.UserGenders, error)
	UserAddGender(gender entities.UserGenders) error
	UserGetAllGender(profileId string) (helperstruct.GenderHelper, error)
	GetPreferenceByProfileId(profileId string) (entities.Preference, error)
	UserAddPreference(entities.Preference) error
	UserEditPreference(entities.Preference) error
	GetUserById(userId string) (entities.User, error)
	UploadProfileImage(Image, ProfileId string) (string, error)
	GetProfilePic(string) (string, error)
	UpdateAge(age int, profileId string) error
	GetAge(profileId string) (int, error)
	FetchUser(profile string) (helperstruct.FetchUser, error)
	FetchPreference(string) (helperstruct.FetchPreference, error)
	FetchInterests(id string) ([]string, error)
	FetchUsers(maxAge, minAge, gender int, id string) ([]helperstruct.Home, error)
	FetchImages(id string) (string, error)

	IsUserExist(id string) (bool, error)
	DecrementLikeCount(userId string) error
	UpdateSubscription(userId string, subscribed bool) error
}
