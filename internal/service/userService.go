package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/akshaybt001/DatingApp_UserService/entities"
	helperstruct "github.com/akshaybt001/DatingApp_UserService/entities/helperStruct"
	"github.com/akshaybt001/DatingApp_UserService/internal/adapters"
	"github.com/akshaybt001/DatingApp_UserService/internal/helper"
	"github.com/akshaybt001/DatingApp_UserService/internal/usecases"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type UserService struct {
	adapters adapters.AdapterInterface
	usecases usecases.Usecases
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapters adapters.AdapterInterface, usecases usecases.Usecases) *UserService {
	return &UserService{
		adapters: adapters,
		usecases: usecases,
	}
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-service:6379",
		Password: "",
		DB:       0,
	})
}

func (user *UserService) UserSignup(ctx context.Context, req *pb.UserSignupRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		logger.Warn("email can't be empty")
		return nil, fmt.Errorf("email can't be empty")
	}
	if req.Name == "" {
		logger.Warn("name cant be empty")
		return nil, fmt.Errorf("name can't be empty")
	}
	if req.Password == "" {
		logger.Warn("password can't be empty")
		return nil, fmt.Errorf("password can't be empty")
	}
	if req.Phone == "" {
		logger.Warn("phone can't be empty")
		return nil, fmt.Errorf("phone can't be empty")
	}
	check1, err := user.adapters.GetUserByEmail(req.Email)
	if err != nil {
		logger.Error("error in fetching email", "email", req.Email)
		return nil, err
	}
	if check1.Name != "" {
		logger.Error("error account already exists with the given email", "email", req.Email)
		return nil, fmt.Errorf("an account already exists with the given email")
	}
	check2, err := user.adapters.GetUserByPhone(req.Phone)
	if err != nil {
		logger.Error("error in fetching userby phone")
		return nil, err
	}
	if check2.Name != "" {
		logger.Error("error account already exists with the given phone", "phone", req.Phone)
		return nil, fmt.Errorf("an account already exist with the given phone number")
	}
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		logger.Error("error in hashing password", "email", req.Email)
		return nil, err
	}
	reqEntity := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}
	res, err := user.adapters.UserSignup(reqEntity)
	if err != nil {
		logger.Error("error in user signup", "email", req.Email)
		return nil, err
	}
	return &pb.UserSignupResponse{
		Id:    res.ID.String(),
		Name:  res.Name,
		Email: res.Email,
		Phone: res.Phone,
	}, nil
}

func (user *UserService) UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		logger.Warn("invalid email", "email-", req.Email)
		return &pb.UserSignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	userData, err := user.adapters.GetUserByEmail(req.Email)
	if err != nil {
		logger.Error("error in fetching userData")
		return &pb.UserSignupResponse{}, err
	}
	if userData.IsBlocked {
		logger.Warn("user have been blocked by the admin", "email", req.Email)
		return &pb.UserSignupResponse{}, fmt.Errorf("you have been blocked by the admin")
	}
	if userData.Email == "" {
		logger.Warn("invalid credentials ", "email", req.Email)
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(userData.Password, req.Password) {
		logger.Error("error in compareing password")
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials please try again")
	}
	return &pb.UserSignupResponse{
		Id:    userData.ID.String(),
		Name:  userData.Name,
		Email: userData.Email,
		Phone: userData.Phone,
	}, nil
}

func (user *UserService) AdminLogin(ctx context.Context, req *pb.LoginRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		logger.Warn("invalid email", "email", req.Email)
		return &pb.UserSignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	adminData, err := user.adapters.GetAdminByEmail(req.Email)
	if err != nil {
		logger.Error("error in fetching admin data")
		return &pb.UserSignupResponse{}, err
	}
	if adminData.Email == "" {
		logger.Warn("invalid credentials")
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(adminData.Password, req.Password) {
		logger.Error("error in compareing password")
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credential")
	}
	return &pb.UserSignupResponse{
		Id:    adminData.ID.String(),
		Name:  adminData.Name,
		Email: adminData.Email,
		Phone: adminData.Phone,
	}, nil
}

func (user *UserService) CreateProfile(ctx context.Context, req *pb.GetUserById) (*pb.NoArg, error) {
	if err := user.adapters.CreateProfile(req.Id); err != nil {
		return &pb.NoArg{}, err
	}
	logger.Info("creating profile for user", "user_id", req.Id)

	return &pb.NoArg{}, nil
}

func (user *UserService) AdminAddInterest(ctx context.Context, req *pb.AddInterestRequest) (*pb.NoArg, error) {
	reqEntity := entities.Interests{
		Interest: req.Interest,
	}
	check1, err := user.adapters.GetInterestByName(req.Interest)
	if err != nil {
		logger.Error("error fectching interest by name", "interest_name", req.Interest, "error", err)
		return nil, err
	}
	if check1.Interest != "" {
		logger.Warn("interest already exist")
		return nil, fmt.Errorf("interest already exist")
	}
	err = user.adapters.AdminAddInterest(reqEntity)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) AdminUpdateInterest(ctx context.Context, req *pb.InterestResponse) (*pb.NoArg, error) {
	reqEntity := entities.Interests{
		Id:       int(req.Id),
		Interest: req.Interest,
	}
	check1, err := user.adapters.GetInterestByName(req.Interest)
	if err != nil {
		logger.Error("error fetching the interest by name")
		return nil, err
	}
	if check1.Interest != "" {
		logger.Warn("interest already exist")
		return nil, fmt.Errorf("interest already exist")
	}
	if err := user.adapters.AdminUpdateInterest(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) AdminUpdateGender(ctx context.Context, req *pb.GenderResponse) (*pb.NoArg, error) {
	reqEntity := entities.Gender{
		Id:   int(req.Id),
		Name: req.Gender,
	}
	check, err := user.adapters.GetGenderByName(req.Gender)
	if err != nil {
		logger.Error("error fetching the gender by name")
		return nil, err
	}
	if check.Name != "" {
		logger.Warn("gender already exist")
		return nil, fmt.Errorf("gender already exist")
	}
	if err := user.adapters.AdminUpdateGender(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) GetAllInterest(e *pb.NoArg, srv pb.UserService_GetAllInterestServer) error {
	interests, err := user.adapters.AdminGetAllInterest()
	if err != nil {
		logger.Error("error in fetching get all interest")
		return err
	}
	for _, interest := range interests {
		res := &pb.InterestResponse{
			Id:       int32(interest.Id),
			Interest: interest.Interest,
		}
		err := srv.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *UserService) GetAllGender(e *pb.NoArg, srv pb.UserService_GetAllGenderServer) error {
	genders, err := user.adapters.AdminGetAllGender()
	if err != nil {
		logger.Error("error in fetching get all gender")
		return err
	}
	for _, gender := range genders {
		res := &pb.GenderResponse{
			Id:     int32(gender.Id),
			Gender: gender.Name,
		}
		err := srv.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *UserService) AddInterestUser(ctx context.Context, req *pb.DeleteInterestRequest) (*pb.NoArg, error) {
	check, err := user.adapters.GetInterestById(int(req.InterestId))
	if err != nil {
		logger.Error("error fectching interest by ID", "interest_id", req.InterestId, "error", err)
		return nil, err
	}
	if check.InterestId == 0 {
		logger.Warn("Invalid interest ID provided", "interest_id", req.InterestId)

		return nil, fmt.Errorf("please enter a valid interest id")
	}
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	loggerctx := logger.With("user_id", req.UserId)
	check1, err := user.adapters.GetUserInterestById(profile, int(req.InterestId))
	if err != nil {
		loggerctx.Error("error fectching user interest", "interest_id", req.InterestId, "error ", err)
		return nil, err
	}
	if check1.InterestId != 0 {
		loggerctx.Warn("interest already added for user", "interest_id", req.InterestId)
		return nil, fmt.Errorf("you already have added this interest please add a new one")

	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	reqEntity := entities.UserInterests{
		ProfileId:  profileId,
		InterestId: int(req.InterestId),
	}
	if err := user.adapters.UserAddInterest(reqEntity); err != nil {
		loggerctx.Error("Error adding user interest", "interest_id", req.InterestId, "error", err)
		return nil, err
	}
	loggerctx.Info("interest added successfully for user", "interest_id", req.InterestId)
	return nil, nil
}

func (user *UserService) AddGenderUser(ctx context.Context, req *pb.UpdateGenderRequest) (*pb.NoArg, error) {
	check, err := user.adapters.GetGenderById(int(req.GenderId))
	if err != nil {
		logger.Error("Error to fectching gender id", "gender_id", req.GenderId, "error", err)
		return nil, err
	}
	if check.GenderId == 0 {
		logger.Error("Error gender_id is not correct ")
		return nil, fmt.Errorf("please enter a valid gender id")
	}
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	loggerctx := logger.With("user_id", req.UserId)
	check1, err := user.adapters.GetUserGenderById(profile, int(req.GenderId))
	if err != nil {
		logger.Error("error fetching the gender_id", "gender_Id", req.GenderId, "error", err)
		return nil, err
	}
	if check1.GenderId != 0 {
		logger.Error("error gender is already added")
		return nil, fmt.Errorf("you already have added this gender ")

	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	reqEntity := entities.UserGenders{
		ProfileId: profileId,
		GenderId:  int(req.GenderId),
	}
	if err := user.adapters.UserAddGender(reqEntity); err != nil {
		loggerctx.Error("Error adding user gender", "gender_id", req.GenderId, "error", err)

		return nil, err
	}
	loggerctx.Info("gender added successfully for user", "gender_id", req.GenderId)
	return nil, nil
}

func (user *UserService) DeleteInterestUser(ctx context.Context, req *pb.DeleteInterestRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	loggerctx := logger.With("user_id", req.UserId)
	reqEntity := entities.UserInterests{
		ProfileId:  profileId,
		InterestId: int(req.InterestId),
	}
	if err := user.adapters.UserDeleteInterest(reqEntity); err != nil {
		loggerctx.Error("Error deleting user interest", "interest_id", req.InterestId, "error", err)
		return nil, err
	}
	loggerctx.Info("gender added successfully for user", "interest_id", req.InterestId)
	return nil, nil
}

func (user *UserService) GetAllInterestsUser(req *pb.GetUserById, srv pb.UserService_GetAllInterestsUserServer) error {
	profileId, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return err
	}
	interests, err := user.adapters.UserGetAllInterest(profileId)
	if err != nil {
		logger.Error("Error in fetching interests")
		return err
	}
	for _, interest := range interests {
		res := &pb.InterestResponse{
			Id:       int32(interest.Id),
			Interest: interest.Interest,
		}

		if err := srv.Send(res); err != nil {
			return err
		}
	}
	logger.Info("fetching interest is successful")
	return nil
}

func (user *UserService) UserAddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	address, err := user.adapters.GetAddressByProfileId(profile)
	if err != nil {
		logger.Error("Error fetching address", "profile_id", profile, "error", err)
		return nil, err
	}
	if address.Country != "" {
		logger.Error("address is already exists")
		return nil, fmt.Errorf("you have already added an address please edit the existing")
	}
	reqEntity := entities.Address{
		Country:   req.Country,
		State:     req.State,
		District:  req.District,
		City:      req.City,
		ProfileId: profileId,
	}
	if err := user.adapters.UserAddAddress(reqEntity); err != nil {
		logger.Error("Error adding user address", "error", err)
		return nil, err
	}
	logger.Info("address added successfully for user")
	return nil, nil
}

func (user *UserService) UserEditAddress(ctx context.Context, req *pb.AddressResponse) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	reqEntity := entities.Address{
		Country:   req.Country,
		State:     req.State,
		District:  req.District,
		City:      req.City,
		ProfileId: profileId,
	}
	if err := user.adapters.UserEditAddress(reqEntity); err != nil {
		logger.Error("Error editing user address", "error", err)
		return nil, err
	}
	logger.Info("address edited successfully for user", "user_id", req.UserId)
	return nil, nil
}

func (user *UserService) UserEditPreference(ctx context.Context, req *pb.PreferenceResponse) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	reqEntity := entities.Preference{
		MinAge:     int(req.Minage),
		MaxAge:     int(req.Maxage),
		GenderId:   int(req.Gender),
		DesireCity: req.Desirecity,
		ProfileId:  profileId,
	}
	if err := user.adapters.UserEditPreference(reqEntity); err != nil {
		logger.Error("Error editing user preference", "error", err)
		return nil, err
	}
	logger.Info("preference edited successfully for user", "user_id", req.UserId)
	return nil, nil

}

func (user *UserService) UserGetAddress(ctx context.Context, req *pb.GetUserById) (*pb.AddressResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return nil, err
	}
	address, err := user.adapters.GetAddressByProfileId(profile)
	if err != nil {
		logger.Error("Error fetching address", "profile_id", profile, "error", err)
		return nil, err
	}
	addressId := ""
	if address.Id != uuid.Nil {
		addressId = address.Id.String()
	}
	res := &pb.AddressResponse{
		Id:       addressId,
		Country:  address.Country,
		State:    address.State,
		District: address.District,
		City:     address.City,
	}
	logger.Info("Address successfully fetched", "user_id", req.Id)
	return res, nil
}
func (user *UserService) GetAllGenderUser(ctx context.Context, req *pb.GetUserById) (*pb.GenderResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return nil, err
	}
	// genders, err := user.adapters.UserGetAllGender(profileId)
	// if err != nil {
	// 	return err
	// }
	// for _, gender := range genders {
	// 	res := &pb.GenderResponse{
	// 		Id:     int32(gender.GenderId),
	// 		Gender: gender.GenderName,
	// 	}
	// 	if err := srv.Send(res); err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	// gender,err:=user.adapters.GetGenderByProfileId(profile)
	// if err!=nil{
	// 	return nil,err
	// }
	genders, err := user.adapters.UserGetAllGender(profile)
	if err != nil {
		logger.Error("error in fetching gender")
		return nil, err
	}
	res := &pb.GenderResponse{
		Id:     int32(genders.GenderId),
		Gender: genders.GenderName,
	}
	logger.Info("successfully fetched gender")
	return res, nil

}

func (user *UserService) GetAllPreference(ctx context.Context, req *pb.GetUserById) (*pb.PreferenceResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return nil, err
	}
	preference, err := user.adapters.GetPreferenceByProfileId(profile)
	if err != nil {
		logger.Error("Error fetching preference", "profile_id", profile, "error", err)
		return nil, err
	}
	preferenceId := ""
	if preference.Id != uuid.Nil {
		preferenceId = preference.Id.String()
	}
	res := &pb.PreferenceResponse{
		Id:         preferenceId,
		Minage:     int32(preference.MinAge),
		Maxage:     int32(preference.MaxAge),
		Gender:     int32(preference.GenderId),
		Desirecity: preference.DesireCity,
	}
	return res, nil
}

func (user *UserService) AdminAddGender(ctx context.Context, req *pb.AddGenderRequest) (*pb.NoArg, error) {
	reqEntity := entities.Gender{
		Name: req.Gender,
	}
	check, err := user.adapters.GetGenderByName(req.Gender)
	if err != nil {
		logger.Error("error fetching gender name")
		return nil, err
	}
	if check.Name != "" {
		logger.Warn("gender already exist")
		return nil, fmt.Errorf("gender already exist")
	}
	err = user.adapters.AdminAddGender(reqEntity)
	if err != nil {
		logger.Error("error in add gender by admin")
		return nil, err
	}
	return nil, nil
}

func (user *UserService) UserAddPreference(ctx context.Context, req *pb.PreferenceRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		logger.Error("Error parsing profile Id, profile_id", profile, "error", err)
		return nil, err
	}
	preference, err := user.adapters.GetPreferenceByProfileId(profile)
	if err != nil {
		logger.Error("Error fetching preference", "profile_id", profile, "error", err)
		return nil, err
	}
	if preference.DesireCity != "" {
		logger.Error("preference is already exists")
		return nil, fmt.Errorf("you have already added a preference please edit the existing")
	}
	reqEntity := entities.Preference{
		MinAge:     int(req.Minage),
		MaxAge:     int(req.Maxage),
		GenderId:   int(req.Gender),
		DesireCity: req.Desirecity,
		ProfileId:  profileId,
	}
	if err := user.adapters.UserAddPreference(reqEntity); err != nil {
		logger.Error("error in adding preference ", "user_id", req.UserId, "error", err)
		return nil, err
	}
	return nil, nil

}

func (user *UserService) GetUser(ctx context.Context, req *pb.GetUserById) (*pb.UserSignupResponse, error) {
	userData, err := user.adapters.GetUserById(req.Id)
	if err != nil {
		logger.Error("error in fetching userid", "used_id", req.Id, "error", err)
		return nil, err
	}
	res := &pb.UserSignupResponse{
		Id:    userData.ID.String(),
		Name:  userData.Name,
		Email: userData.Email,
		Phone: userData.Phone,
	}
	return res, nil
}

func (user *UserService) UserUploadProfileImage(ctx context.Context, req *pb.UserImageRequest) (*pb.UserImageResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}
	url, err := user.usecases.UploadImage(req, profile)
	if err != nil {
		logger.Error("error in uploadimage on usecase")
		return nil, err
	}
	res := &pb.UserImageResponse{
		Url: url,
	}
	return res, nil
}

func (user *UserService) UserGetProfilePic(ctx context.Context, req *pb.GetUserById) (*pb.UserImageResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return nil, err
	}
	image, err := user.adapters.GetProfilePic(profile)
	if err != nil {
		logger.Error("error in fetching profile pic")
		return nil, err
	}
	return &pb.UserImageResponse{
		Url: image,
	}, nil
}

func (user *UserService) UserAddAge(ctx context.Context, req *pb.UserAgeRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.UserId, "error", err)
		return nil, err
	}

	layout := "2006-01-02T15:04:05.999999Z"
	dob, err := time.Parse(layout, req.Dob)
	if err != nil {
		logger.Warn("invalid time format")
		return &pb.NoArg{}, fmt.Errorf("please provide time in appropriate format")
	}
	age := helper.CalculateAge(dob)

	if err := user.adapters.UpdateAge(age, profile); err != nil {
		logger.Error("error in setting age ")
		return nil, err
	}
	return nil, nil

}

func (user *UserService) UserGetAge(ctx context.Context, req *pb.GetUserById) (*pb.UserAgeResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return nil, err
	}
	age, err := user.adapters.GetAge(profile)
	if err != nil {
		logger.Error("error in fetching age", "user_id", req.Id, "error", err)
		return nil, err
	}
	return &pb.UserAgeResponse{
		Age: int32(age),
	}, nil
}

func (user *UserService) HomePage(ctx context.Context, req *pb.GetUserById) (*pb.HomeResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
		return nil, err
	}
	preference, err := user.adapters.FetchPreference(profile)
	if err != nil {
		logger.Error("error fetching preference by userId", "user_id", req.Id, "error", err)
		return nil, err
	}
	userData, err := user.adapters.FetchUser(profile)
	if err != nil {
		logger.Error("error fetching userData by userId", "user_id", req.Id, "error", err)
		return nil, err
	}
	interestData, err := user.adapters.FetchInterests(profile)
	if err != nil {
		logger.Error("error fetching interests by userId", "user_id", req.Id, "error", err)
		return nil, err
	}
	end := len(interestData) - 1
	users, err := user.adapters.FetchUsers(preference.MinAge, preference.MaxAge, preference.Gender, profile)
	if err != nil {
		logger.Error("error to fetching users based on preferences", "user_id", req.Id, "error", err)
		return nil, err
	}

	displayedUserIds, err := getDisplayedUserIds(profile)
	if err != nil {
		return nil, err
	}
	fmt.Println("display :", displayedUserIds)

	scores := []float64{}
	matchUsers := []helperstruct.Home{}
	for _, u := range users {
		userProfile, err := user.adapters.GetProfileIdByUserId(u.Id)
		if err != nil {
			logger.Error("error fetching profile ID by user ID", "user_id", req.Id, "error", err)
			return nil, err
		}
		image, err := user.adapters.FetchImages(userProfile)
		if err != nil {
			logger.Error("error fetching images ", "user_id", u.Id)
			return nil, err
		}
		u.Images = image

		interests, err := user.adapters.FetchInterests(userProfile)
		if err != nil {
			logger.Error("error fetching images", "user_id", u.Id)
			return nil, err
		}
		userPreference, err := user.adapters.FetchPreference(userProfile)
		if err != nil {
			logger.Error("error fetching preference", "user_id", u.Id)
			return nil, err
		}
		if userPreference.DesireCity != preference.DesireCity {
			continue
		}

		u.Interests = interests
		interestScore := 0
		for _, interest := range interests {
			search := helper.SearchForInterest(interestData, interest, 0, end)
			if !search {
				continue
			}
			interestScore++
		}
		fmt.Println("akshay :", u.Name)
		if !displayedUserIds[u.Id] {
			AgeScore := helper.Abs(u.Age - userData.Age)
			score := float64(AgeScore) + 2*float64(interestScore)
			scores = append(scores, score)
			matchUsers = append(matchUsers, u)

			displayedUserIds[u.Id] = true
		}

	}
	if len(matchUsers) == 0 {
		logger.Error("there is no new recommendations")
		return nil, fmt.Errorf("no new recommendations available")
	}

	helper.QuickSort(scores, matchUsers, 0, len(scores)-1)

	var homeResponse *pb.HomeResponse
	id := matchUsers[0].Id

	seen := map[string]bool{id: true}

	err = updateDisplayedUserIds(profile, seen)
	if err != nil {
		return nil, fmt.Errorf("here is the problem in update displayUserId")
	}

	homeResponse = &pb.HomeResponse{
		Id:        matchUsers[0].Id,
		Name:      matchUsers[0].Name,
		Age:       int32(matchUsers[0].Age),
		Gender:    matchUsers[0].Gender,
		City:      matchUsers[0].City,
		Country:   matchUsers[0].Country,
		Image:     matchUsers[0].Images,
		Interests: matchUsers[0].Interests,
	}

	fmt.Println("home :", homeResponse)

	return homeResponse, nil
}

func getDisplayedUserIds(userID string) (map[string]bool, error) {
	displayedUserIdsKey := fmt.Sprintf("displayed_user_ids:%s", userID)
	displayedUserIds := make(map[string]bool)

	userIdsInRedis, err := redisClient.SMembers(displayedUserIdsKey).Result()
	if err != nil {
		return nil, err
	}

	for _, userIdStr := range userIdsInRedis {
		displayedUserIds[userIdStr] = true
	}

	return displayedUserIds, nil
}

func updateDisplayedUserIds(userID string, displayedUserIds map[string]bool) error {
	displayedUserIdsKey := fmt.Sprintf("displayed_user_ids:%s", userID)
	userIds := make([]interface{}, 0, len(displayedUserIds))

	for userId := range displayedUserIds {
		userIds = append(userIds, userId)
	}

	_, err := redisClient.SAdd(displayedUserIdsKey, userIds...).Result()
	return err
}

func (user *UserService) GetUserData(ctx context.Context, req *pb.GetUserById) (*pb.UserDataResponse, error) {
	userData, err := user.adapters.GetUserById(req.Id)
	if err != nil {
		logger.Error("error fetching userData", "user_id", req.Id, "error", err)
		return nil, err
	}
	res := &pb.UserDataResponse{
		Id:           userData.ID.String(),
		Name:         userData.Name,
		Email:        userData.Email,
		Phone:        userData.Phone,
		IsBlocked:    userData.IsBlocked,
		LikeCount:    int32(userData.LikeCount),
		IsSubscribed: userData.IsSubscribed,
	}
	return res, nil
}

func (user *UserService) DecrementLikeCount(ctx context.Context, req *pb.GetUserById) (*pb.NoArg, error) {
	err := user.adapters.DecrementLikeCount(req.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) UpdateSubscription(ctx context.Context, req *pb.UpdateSubscriptionRequest) (*pb.NoArg, error) {
	if err := user.adapters.UpdateSubscription(req.UserId, req.Subscription); err != nil {
		return nil, err
	}
	return nil, nil
}
