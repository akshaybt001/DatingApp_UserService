package usecases

import "github.com/akshaybt001/DatingApp_proto_files/pb"

type Usecases interface {
	UploadImage(*pb.UserImageRequest,string)(string,error)
	// UpdateDisplayedUserIds(userID string, displayedUserIds map[string]bool) error
	// GetDisplayedUserIds(userID string) (map[string]bool, error) 
}