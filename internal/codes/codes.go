package codes

import (
	"errors"
)

var (
	// env
	ErrEnvNoAcsess          = errors.New("no access to object")
	ErrEnvNotStructure      = errors.New("object is not a structure")
	ErrEnvFieldNotSet       = errors.New("field cannot be set")
	ErrEnvFieldNotSupported = errors.New("field type not supported")
	ErrEnvFieldParseUint    = errors.New("field uint parse fail")
	ErrEnvFieldParseInt     = errors.New("field int parse fail")
	//flags
	ErrFlgNoAcsess          = errors.New("no access to object")
	ErrFlgNotStructure      = errors.New("object is not a structure")
	ErrFlgFieldNotSet       = errors.New("field cannot be set")
	ErrFlgFieldNotSupported = errors.New("field type not supported")
	ErrFlgParseFlag         = errors.New("flag parse fail")
	// memory
	ErrRepParseFloat         = errors.New("value float parse fail")
	ErrRepParseInt           = errors.New("value int parse fail")
	ErrRepNotFound           = errors.New("not found")
	ErrRepMetricNotSupported = errors.New("this type of metrics is not supported")
	// get
	ErrGetValReqType = errors.New("bad request string - type")
	ErrGetValReqName = errors.New("bad request string - name")
	// post
	ErrPostValReqType  = errors.New("bad request string - type")
	ErrPostValReqName  = errors.New("bad request string - name")
	ErrPostValReqValue = errors.New("bad request string - value")
	// get list
	ErrGetLstParse   = errors.New("parse template faild")
	ErrGetLstReqType = errors.New("bad request string - type")
	ErrGetLstReqName = errors.New("bad request string - name")
	// post json
	ErrPostNotJSON   = errors.New("bad request - request not JSON")
	ErrPostBadBody   = errors.New("bad request - bad body")
	ErrPostMarshal   = errors.New("error of marshaling")
	ErrPostUnmarshal = errors.New("error of unmarshaling")
	// compressor
	ErrCompressor = errors.New("bad compressor init")
	// decompressor
	ErrDecompressor = errors.New("bad decompressor init")
	// hasher
	ErrHashNotBody    = errors.New("bad request - body is nil")
	ErrHashNilHeader  = errors.New("bad request - hash in header is nil")
	ErrHashNotCorrect = errors.New("bad request - hash isn`t correct")
)
