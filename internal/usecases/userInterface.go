package usecases

import "github.com/akshaybt001/DatingApp_proto_files/pb"

type Usecases interface {
	UploadImage(*pb.UserImageRequest,string)(string,error)
}