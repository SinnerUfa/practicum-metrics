package codes

import (
	"errors"
)

var (
	// env
	ErrEnvNoAcsess          = errors.New("env: no access to object")
	ErrEnvNotStructure      = errors.New("env: object is not a structure")
	ErrEnvFieldNotSet       = errors.New("env: field cannot be set")
	ErrEnvFieldNotSupported = errors.New("env: field type not supported")
	ErrEnvFieldParseUint    = errors.New("env: field uint parse fail")
	ErrEnvFieldParseInt     = errors.New("env: field int parse fail")
	//flags
	ErrFlgNoAcsess          = errors.New("flag: no access to object")
	ErrFlgNotStructure      = errors.New("flag: object is not a structure")
	ErrFlgFieldNotSet       = errors.New("flag: field cannot be set")
	ErrFlgFieldNotSupported = errors.New("flag: field type not supported")
	ErrFlgParseFlag         = errors.New("flag: flag parse fail")
	// memory
	ErrRepParseFloat         = errors.New("rep: value float parse fail")
	ErrRepParseInt           = errors.New("rep: value int parse fail")
	ErrRepNotFound           = errors.New("rep: not found")
	ErrRepMetricNotSupported = errors.New("rep: this type of metrics is not supported")
	// get
	ErrGetValReqType = errors.New("get value: bad request string - type")
	ErrGetValReqName = errors.New("get value: bad request string - name")
	// post
	ErrPostValReqType  = errors.New("post value: bad request string - type")
	ErrPostValReqName  = errors.New("post value: bad request string - name")
	ErrPostValReqValue = errors.New("post value: bad request string - value")
	// get list
	ErrGetLstParse   = errors.New("get list: parse template faild")
	ErrGetLstReqType = errors.New("get list: bad request string - type")
	ErrGetLstReqName = errors.New("get list: bad request string - name")
)