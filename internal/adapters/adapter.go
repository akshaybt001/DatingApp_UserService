package adapters

import (
	"github.com/akshaybt001/DatingApp_UserService/entities"
	helperstruct "github.com/akshaybt001/DatingApp_UserService/entities/helperStruct"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func NewUserAdapter(db *gorm.DB) *UserAdapter {
	return &UserAdapter{
		DB: db,
	}
}

func (user *UserAdapter) UserSignup(userData entities.User) (entities.User, error) {
	var res entities.User
	id := uuid.New()
	insertQuery := `INSERT INTO users (id,name,email,password,phone,created_at) VALUES ($1,$2,$3,$4,$5,NOW()) RETURNING *`
	if err := user.DB.Raw(insertQuery, id, userData.Name, userData.Email, userData.Password, userData.Phone).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetUserByEmail(email string) (entities.User, error) {
	var res entities.User
	selectQuery := `SELECT * FROM users WHERE email=?`
	if err := user.DB.Raw(selectQuery, email).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetUserByPhone(phone string) (entities.User, error) {
	var res entities.User
	selectQuery := `SELECT * FROM users WHERE phone=?`
	if err := user.DB.Raw(selectQuery, phone).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetAdminByEmail(email string) (entities.Admin, error) {
	var res entities.Admin
	selectQuery := `SELECT * FROM admins WHERE email=?`
	if err := user.DB.Raw(selectQuery, email).Scan(&res).Error; err != nil {
		return entities.Admin{}, err
	}
	return res, nil
}

func (user *UserAdapter) CreateProfile(userID string) error {
	profileId := uuid.New()
	insertProfile := `INSERT INTO profiles (id,user_id) VALUES ($1,$2)`
	if err := user.DB.Exec(insertProfile, profileId, userID).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) GetProfileIdByUserId(userId string) (string, error) {
	var profileId string
	selectProfile := `SELECT id FROM profiles WHERE user_id=?`
	if err := user.DB.Raw(selectProfile, userId).Scan(&profileId).Error; err != nil {
		return "", err
	}
	return profileId, nil
}

func (user *UserAdapter) AdminAddInterest(interest entities.Interests) error {
	var id int
	selectMaxId := `SELECT COALESCE(MAX(id),0) FROM interests`
	if err := user.DB.Raw(selectMaxId).Scan(&id).Error; err != nil {
		return err
	}
	insertInterest := `INSERT INTO interests (id,interest) VALUES ($1,$2)`
	if err := user.DB.Exec(insertInterest, id+1, interest.Interest).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) AdminAddGender(gender entities.Gender) error {
	var id int
	selectMaxId := `SELECT COALESCE(MAX(id),0) FROM genders`
	if err := user.DB.Raw(selectMaxId).Scan(&id).Error; err != nil {
		return err
	}
	insertGender := `INSERT INTO genders (id,name) VALUES ($1,$2)`
	if err := user.DB.Exec(insertGender, id+1, gender.Name).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) GetInterestByName(interest string) (entities.Interests, error) {
	var res entities.Interests
	selectQuery := `SELECT * FROM interests WHERE interest=?`
	if err := user.DB.Raw(selectQuery, interest).Scan(&res).Error; err != nil {
		return entities.Interests{}, err
	}
	return res, nil
}


func (user *UserAdapter) GetGenderByName(gender string) (entities.Gender, error) {
	var res entities.Gender
	selectQuery := `SELECT * FROM genders WHERE name=?`
	if err := user.DB.Raw(selectQuery, gender).Scan(&res).Error; err != nil {
		return entities.Gender{}, err
	}
	return res, nil
}

func (user *UserAdapter) AdminUpdateInterest(interest entities.Interests) error {
	updateInterest := `UPDATE interests SET interest=$1 WHERE id=$2`
	if err := user.DB.Exec(updateInterest, interest.Interest, interest.Id).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) AdminUpdateGender(gender entities.Gender) error {
	updateGender := `UPDATE genders SET name=$1 WHERE id=$2`
	if err := user.DB.Exec(updateGender, gender.Name, gender.Id).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) AdminGetAllInterest() ([]entities.Interests, error) {
	var res []entities.Interests
	selectInterest := `SELECT * FROM interests`
	if err := user.DB.Raw(selectInterest).Scan(&res).Error; err != nil {
		return []entities.Interests{}, err
	}
	return res, nil
}

func (user *UserAdapter) AdminGetAllGender() ([]entities.Gender, error) {
	var res []entities.Gender
	selectGender := `SELECT * FROM genders`
	if err := user.DB.Raw(selectGender).Scan(&res).Error; err != nil {
		return []entities.Gender{}, err
	}
	return res, nil
}

func (user *UserAdapter) UserAddInterest(interests entities.UserInterests) error {
	var id int
	selectMaxId := `SELECT COALESCE(MAX(id),0) FROM user_interests`
	if err := user.DB.Raw(selectMaxId).Scan(&id).Error; err != nil {
		return err
	}
	insertInterestQuery := `INSERT INTO user_interests(id,interest_id,profile_id) VALUES ($1,$2,$3)`
	if err := user.DB.Exec(insertInterestQuery, id+1, interests.InterestId, interests.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

// UserAddGender implements AdapterInterface.
func (user *UserAdapter) UserAddGender(gender entities.UserGenders) error {
	var id int
	selectMaxId := `SELECT COALESCE(MAX(id),0) FROM user_genders`
	if err := user.DB.Raw(selectMaxId).Scan(&id).Error; err != nil {
		return err
	}
	insertInterestQuery := `INSERT INTO user_genders(id,gender_id,profile_id) VALUES ($1,$2,$3)`
	if err := user.DB.Exec(insertInterestQuery, id+1, gender.GenderId, gender.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) UserDeleteInterest(interest entities.UserInterests) error {
	deleteInterestQuery := `DELETE FROM user_interests WHERE interest_id=$1 AND profile_id=$2`
	if err := user.DB.Exec(deleteInterestQuery, interest.InterestId, interest.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) UserGetAllInterest(profileId string) ([]entities.Interests, error) {
	var res []entities.Interests
	selectInterestQueryUser := `SELECT i.id ,i.interest FROM interests i JOIN user_interests u ON u.interest_id=i.id WHERE profile_id=$1`
	if err := user.DB.Raw(selectInterestQueryUser, profileId).Scan(&res).Error; err != nil {
		return []entities.Interests{}, err
	}
	return res, nil
}

func (user *UserAdapter) UserGetAllGender(profileId string) (helperstruct.GenderHelper, error) {
	var res helperstruct.GenderHelper
	selectQueryUser := `SELECT g.id AS gender_id,g.name AS gender_name FROM genders g JOIN user_genders u ON u.gender_id=g.id WHERE profile_id=$1`
	if err := user.DB.Raw(selectQueryUser, profileId).Scan(&res).Error; err != nil {
		return helperstruct.GenderHelper{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetInterestById(id int) (helperstruct.InterestHelper, error) {
	selectInterestQuery := `SELECT id AS interest_id , interest AS interest_name FROM interests WHERE id=?`
	var res helperstruct.InterestHelper
	if err := user.DB.Raw(selectInterestQuery, id).Scan(&res).Error; err != nil {
		return helperstruct.InterestHelper{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetGenderById(id int) (helperstruct.GenderHelper, error) {
	selectGenderQuery := `SELECT id AS gender_id, name AS gender_name FROM genders WHERE id=?`
	var res helperstruct.GenderHelper
	if err := user.DB.Raw(selectGenderQuery, id).Scan(&res).Error; err != nil {
		return helperstruct.GenderHelper{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetUserInterestById(profileId string, interesetId int) (entities.UserInterests, error) {
	var res entities.UserInterests
	selectQuery := `SELECT * FROM user_interests WHERE profile_id=$1 AND interest_id=$2`
	if err := user.DB.Raw(selectQuery, profileId, interesetId).Scan(&res).Error; err != nil {
		return entities.UserInterests{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetUserGenderById(profileId string, genderId int) (entities.UserGenders, error) {
	var res entities.UserGenders
	selectQuery := `SELECT * FROM user_genders WHERE profile_id=$1 AND gender_id=$2`
	if err := user.DB.Raw(selectQuery, profileId, genderId).Scan(&res).Error; err != nil {
		return entities.UserGenders{}, err
	}
	return res, nil
}

func (user *UserAdapter) UserAddAddress(req entities.Address) error {
	id := uuid.New()
	insertQuery := `INSERT INTO addresses (id,country,state,district,city,profile_id) VALUES ($1,$2,$3,$4,$5,$6)`
	if err := user.DB.Exec(insertQuery, id, req.Country, req.State, req.District, req.City, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) UserAddPreference(req entities.Preference) error {
	id := uuid.New()
	insertQuery := `INSERT INTO preferences (id,min_age,max_age,gender_id,desire_city,profile_id) VALUES ($1,$2,$3,$4,$5,$6)`
	if err := user.DB.Exec(insertQuery, id, req.MinAge, req.MaxAge, req.GenderId, req.DesireCity, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) UserEditAddress(req entities.Address) error {
	updateQuery := `UPDATE addresses SET country=$1,state=$2,district=$3,city=$4 WHERE profile_id=$5`
	if err := user.DB.Exec(updateQuery, req.Country, req.State, req.District, req.City, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

// UserEditPreference implements AdapterInterface.
func (user *UserAdapter) UserEditPreference(req entities.Preference) error {
	updateQuery := `UPDATE preferences SET min_age=$1,max_age=$2,gender_id=$3,desire_city=$4 WHERE profile_id=$5`
	if err := user.DB.Exec(updateQuery, req.MinAge, req.MaxAge, req.GenderId, req.DesireCity, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) GetAddressByProfileId(id string) (entities.Address, error) {
	var res entities.Address
	selectQuery := `SELECT * FROM addresses WHERE profile_id=?`
	if err := user.DB.Raw(selectQuery, id).Scan(&res).Error; err != nil {
		return entities.Address{}, err
	}
	return res, nil
}

func (user *UserAdapter) GetGenderByProfileId(id string) (entities.UserGenders, error) {
	var res entities.UserGenders
	selectQuery := `SELECT * FROM user_genders WHERE profile_id=?`
	if err := user.DB.Raw(selectQuery, id).Scan(&res).Error; err != nil {
		return entities.UserGenders{}, err
	}
	return res, nil
}

// GetPreferenceByProfileId implements AdapterInterface.
func (user *UserAdapter) GetPreferenceByProfileId(profileId string) (entities.Preference, error) {
	var res entities.Preference
	selectQuery := `SELECT * FROM preferences WHERE profile_id=?`
	if err := user.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return entities.Preference{}, err
	}
	return res, nil
}

// GetUserById implements AdapterInterface.
func (user *UserAdapter) GetUserById(userId string) (entities.User, error) {
	selectUserByIdQuery := `SELECT * FROM users WHERE id=?`
	var res entities.User
	if err := user.DB.Raw(selectUserByIdQuery, userId).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}

func (user *UserAdapter) UploadProfileImage(image, profileId string) (string, error) {
	var res string
	id := uuid.New()
	insertImageQuery := `UPDATE profiles SET image=$1 WHERE id=$2 RETURNING image`
	if err := user.DB.Raw(insertImageQuery, image, profileId).Scan(&res).Error; err != nil {
		return "", err
	}
	insertImageDb := `INSERT INTO images (id,profile_id,file_name) VALUES ($1,$2,$3) `
	if err := user.DB.Exec(insertImageDb, id, profileId, image).Error; err != nil {
		return "", err
	}
	return res, nil
}

func (user *UserAdapter) GetProfilePic(profileId string) (string, error) {
	var res string
	selectQuery := `SELECT image from profiles WHERE id=$1 AND image IS NOT NULL`
	if err := user.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}

func (user *UserAdapter) UpdateAge(age int, profileId string) error {
	var res int
	insertImageQuery := `UPDATE profiles SET age=$1 WHERE id=$2 RETURNING age`
	if err := user.DB.Raw(insertImageQuery, age, profileId).Scan(&res).Error; err != nil {
		return err
	}
	return nil
}

func (user *UserAdapter) GetAge(profileId string) (int, error) {
	var res int
	selectQuery := `SELECT age from profiles WHERE id=$1 AND age IS NOT NULL`
	if err := user.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return 0, err
	}
	return res, nil
}

func (user *UserAdapter) FetchUser(profileId string) (helperstruct.FetchUser, error) {
	var res helperstruct.FetchUser
	selectQuery := `SELECT age from profiles WHERE id=$1`
	if err := user.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return helperstruct.FetchUser{}, err
	}
	return res, nil
}

func (user *UserAdapter) FetchPreference(profileId string) (helperstruct.FetchPreference, error) {
	var res helperstruct.FetchPreference
	selectQuery := `SELECT min_age,max_age,gender_id AS gender,desire_city FROM preferences WHERE profile_id=?`
	if err := user.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return helperstruct.FetchPreference{}, err
	}
	return res, nil
}

func (user *UserAdapter) FetchInterests(id string) ([]string, error) {
	var interests []string
	selectQuery := `SELECT i.interest FROM interests i JOIN user_interests u ON u.interest_id=i.id WHERE profile_id=$1`
	if err := user.DB.Raw(selectQuery, id).Scan(&interests).Error; err != nil {
		return nil, err
	}
	return interests, nil
}

func (user *UserAdapter) FetchUsers(maxAge, minAge, gender int, id string) ([]helperstruct.Home, error) {
	var users []helperstruct.Home
	selectQuery := `SELECT u.id ,u.name , p.age , g.name as gender, a.city , a.country ,p.image  FROM users u JOIN profiles p ON u.id=p.user_id JOIN user_genders ug ON p.id=ug.profile_id JOIN genders g ON g.id=ug.gender_id JOIN addresses a ON p.id=a.profile_id WHERE p.age>? AND p.age<? AND g.id=? AND p.id!=?`
	if err := user.DB.Raw(selectQuery, maxAge, minAge, gender, id).Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (user *UserAdapter) FetchImages(id string) ([]string, error) {
	var images []string
	selectQuery := `SELECT file_name FROM images i JOIN profiles p ON i.profile_id=p.id WHERE profile_id=?`
	if err := user.DB.Raw(selectQuery, id).Scan(&images).Error; err != nil {
		return []string{}, err
	}
	return images, nil
}

func (user *UserAdapter) IsUserExist(id string) (bool, error) {
	var count int
	if err := user.DB.Raw(`SELECT COUNT(*) FROM users WHERE id=?`, id).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (user *UserAdapter) DecrementLikeCount(userId string) error {
	query := "UPDATE users SET like_count = like_count - 1 WHERE id = ?"
	if err := user.DB.Exec(query, userId).Error; err != nil {
		return err
	}

	return nil
}

func (user *UserAdapter) UpdateSubscription(userId string, subscribed bool) error {
	updateQuery := `UPDATE users SET is_subscribed=$1 WHERE id=$2`
	if err := user.DB.Exec(updateQuery, subscribed, userId).Error; err != nil {
		return err
	}
	return nil
}
