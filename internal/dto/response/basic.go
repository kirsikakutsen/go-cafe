package response

import "time"

type basicResponseDto struct {
	Data any `json:"data"`
	Time time.Time `json:"time"`
}

func NewBasicSuccessDto(data any) basicResponseDto {
	return basicResponseDto{
		Data: data,
		Time: time.Now(),
	}
}

func NewBasicErrorDto(err error) basicResponseDto {
	return basicResponseDto{
		Data: err.Error(),
		Time: time.Now(),
	}
}