package service

import (
	"context"
	"fmt"

	"github.com/akshaybt001/DatingApp_UserService/entities"
	"github.com/akshaybt001/DatingApp_UserService/internal/adapters"
	"github.com/akshaybt001/DatingApp_UserService/internal/helper"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/google/uuid"
)

type UserService struct {
	adapters adapters.AdapterInterface
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapters adapters.AdapterInterface) *UserService {
	return &UserService{
		adapters: adapters,
	}
}

func (user *UserService) UserSignup(ctx context.Context, req *pb.UserSignupRequest) (*pb.UserSignupResponse, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("email can't be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name can't be empty")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password can't be empty")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("phone can't be empty")
	}
	check1, err := user.adapters.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if check1.Name != "" {
		return nil, fmt.Errorf("an account already exists with the given email")
	}
	check2, err := user.adapters.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if check2.Name != "" {
		return nil, fmt.Errorf("an account already exist with the given phone number")
	}
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
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
		return &pb.UserSignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	userData, err := user.adapters.GetUserByEmail(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	if userData.IsBlocked {
		return &pb.UserSignupResponse{}, fmt.Errorf("you have been blocked by the admin")
	}
	if userData.Email == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(userData.Password, req.Password) {
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
		return &pb.UserSignupResponse{}, fmt.Errorf("please enter a valid email")
	}
	adminData, err := user.adapters.GetAdminByEmail(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	if adminData.Email == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CompareHashedPassword(adminData.Password, req.Password) {
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
	return &pb.NoArg{}, nil
}

func (user *UserService) AdminAddInterest(ctx context.Context, req *pb.AddInterestRequest) (*pb.NoArg, error) {
	reqEntity := entities.Interests{
		Interest: req.Interest,
	}
	check1, err := user.adapters.GetInterestByName(req.Interest)
	if err != nil {
		return nil, err
	}
	if check1.Interest != "" {
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
		return nil, err
	}
	if check1.Interest != "" {
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
		return nil, err
	}
	if check.Name != "" {
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
		return nil, err
	}
	if check.InterestId == 0 {
		return nil, fmt.Errorf("please enter a valid interest id")
	}
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	check1, err := user.adapters.GetUserInterestById(profile, int(req.InterestId))
	if err != nil {
		return nil, err
	}
	if check1.InterestId != 0 {
		return nil, fmt.Errorf("you already have added this interest please add a new one")

	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.UserInterests{
		ProfileId:  profileId,
		InterestId: int(req.InterestId),
	}
	if err := user.adapters.UserAddInterest(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) AddGenderUser(ctx context.Context, req *pb.UpdateGenderRequest) (*pb.NoArg, error) {
	check, err := user.adapters.GetGenderById(int(req.GenderId))
	if err != nil {
		return nil, err
	}
	if check.GenderId == 0 {
		return nil, fmt.Errorf("please enter a valid gender id")
	}
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	check1, err := user.adapters.GetUserGenderById(profile, int(req.GenderId))
	if err != nil {
		return nil, err
	}
	if check1.GenderId != 0 {
		return nil, fmt.Errorf("you already have added this gender ")

	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.UserGenders{
		ProfileId: profileId,
		GenderId:  int(req.GenderId),
	}
	if err := user.adapters.UserAddGender(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) DeleteSkillUser(ctx context.Context, req *pb.DeleteInterestRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.UserInterests{
		ProfileId:  profileId,
		InterestId: int(req.InterestId),
	}
	if err := user.adapters.UserDeleteInterest(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) GetAllInterestsUser(req *pb.GetUserById, srv pb.UserService_GetAllInterestsUserServer) error {
	profileId, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return err
	}
	interests, err := user.adapters.UserGetAllInterest(profileId)
	if err != nil {
		return err
	}
	for _, interest := range interests {
		res := &pb.InterestResponse{
			Id:       int32(interest.InterestId),
			Interest: interest.InterestName,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (user *UserService) GetAllGenderUser(req *pb.GetUserById, srv pb.UserService_GetAllGenderUserServer) error {
	profileId, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return err
	}
	genders, err := user.adapters.UserGetAllGender(profileId)
	if err != nil {
		return err
	}
	for _, gender := range genders {
		res := &pb.GenderResponse{
			Id:     int32(gender.GenderId),
			Gender: gender.GenderName,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (user *UserService) UserAddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	address, err := user.adapters.GetAddressByProfileId(profile)
	if err != nil {
		return nil, err
	}
	if address.Country != "" {
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
		return nil, err
	}
	return nil, nil
}

func (user *UserService) UserEditAddress(ctx context.Context, req *pb.AddressResponse) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
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
		return nil, err
	}
	return nil, nil
}

func (user *UserService) UserEditPreference(ctx context.Context, req *pb.PreferenceResponse) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
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
		return nil, err
	}
	return nil, nil

}

func (user *UserService) UserGetAddress(ctx context.Context, req *pb.GetUserById) (*pb.AddressResponse, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.Id)
	if err != nil {
		return nil, err
	}
	address, err := user.adapters.GetAddressByProfileId(profile)
	if err != nil {
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
	return res, nil
}

// func (user *UserService) GetAllPreference(ctx context.Context,req *pb.GetUserById) (*pb.PreferenceResponse,error){
// 	profile,err:=user.adapters.GetProfileIdByUserId(req.Id)
// 	if err!=nil{
// 		return nil,err
// 	}
// 	preference,err:=user.adapters.GetPreferenceByProfileId(profile)
// 	if err!=nil{
// 		return nil,err
// 	}
// 	preferenceId:=""
// 	if preference.Id!=uuid.Nil{
// 		preferenceId=preference.Id.String()
// 	}
// 	res:=&pb.PreferenceResponse{
// 		Id: preferenceId,
// 		Minage: int32(preference.MinAge),
// 		Maxage: int32(preference.MaxAge),
// 		Gender: int32(preference.GenderId),
// 		Desirecity: preference.DesireCity,

// 	}
// 	return res,nil
// }

func (user *UserService) AdminAddGender(ctx context.Context, req *pb.AddGenderRequest) (*pb.NoArg, error) {
	reqEntity := entities.Gender{
		Name: req.Gender,
	}
	check, err := user.adapters.GetGenderByName(req.Gender)
	if err != nil {
		return nil, err
	}
	if check.Name != "" {
		return nil, fmt.Errorf("gender already exist")
	}
	err = user.adapters.AdminAddGender(reqEntity)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (user *UserService) UserAddPreference(ctx context.Context, req *pb.PreferenceRequest) (*pb.NoArg, error) {
	profile, err := user.adapters.GetProfileIdByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	profileId, err := uuid.Parse(profile)
	if err != nil {
		return nil, err
	}
	preference, err := user.adapters.GetPreferenceByProfileId(profile)
	if err != nil {
		return nil, err
	}
	if preference.DesireCity != "" {
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
		return nil, err
	}
	return nil, nil

}
