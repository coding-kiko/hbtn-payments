package errors

/*
{
	"error": {
		"code": 400
		"description": "bad request"
	}
}
*/

type Response struct {
	Error Error `json:"error"`
}

type Error struct {
	Description string `json:"description"`
	Code        int    `json:"code"`
}

// create http error response based on the type of error and send response
func CreateResponse(e error) (int, Response) {
	var resp Response
	var err Error

	switch e.(type) {
	case *NotFound:
		err.Description = e.Error()
		err.Code = 404
		resp.Error = err
	case *FileError:
		err.Description = e.Error()
		err.Code = 400
		resp.Error = err
	case *MethodNotAllowed:
		err.Description = e.Error()
		err.Code = 405
		resp.Error = err
	case *BadRequest:
		err.Description = e.Error()
		err.Code = 400
		resp.Error = err
	case *InternalServer:
		err.Description = e.Error()
		err.Code = 500
		resp.Error = err
	default:
		err.Description = e.Error()
		err.Code = 500
		resp.Error = err
	}
	return resp.Error.Code, resp
}
