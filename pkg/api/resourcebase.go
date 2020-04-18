package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/felipecurvelo/weather-reporting-api/pkg/authorizer"
	"github.com/felipecurvelo/weather-reporting-api/pkg/internalerror"
)

type ResourceBase struct {
}

func (b *ResourceBase) ParseFromBody(r *http.Request, requestModel interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, requestModel)
	if err != nil {
		return fmt.Errorf("Invalid request body (%s)", err.Error())
	}

	return nil
}

func (b *ResourceBase) ValidateAuthToken(ctx context.Context, r *http.Request) error {
	auth := authorizer.FromContext(ctx)
	if auth == nil {
		return internalerror.New("Internal Server Error")
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		return internalerror.New("Empty Token")
	}
	if !auth.ValidateToken(token) {
		return internalerror.New("Invalid Token")
	}

	return nil
}

func (r *ResourceBase) SetResponse(status int, response interface{}, w http.ResponseWriter) {
	b := response

	_, ok := response.([]byte)
	if !ok {
		b, _ = json.Marshal(response)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b.([]byte))
}
