package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	urlpkg "net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

//////////////
// GIN DATA //
//////////////

type GinDataInterface interface {

	// Fetch the test context
	GetCtx() *gin.Context

	// Fetch the test engine
	GetEngine() *gin.Engine

	// Fetch the test response recorder
	GetRecorder() *httptest.ResponseRecorder

	/*
		Initialize the context, engine, and recorder. Is automatically
		called by the tester if using the NewGinTester() function.
	*/
	PrepareForTest()
}

// Simple GIN context intialization
type GinData struct {
	Ctx      *gin.Context
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder
}

func (d *GinData) GetCtx() *gin.Context {
	return d.Ctx
}
func (d *GinData) GetEngine() *gin.Engine {
	return d.Engine
}
func (d *GinData) GetRecorder() *httptest.ResponseRecorder {
	return d.Recorder
}
func (d *GinData) PrepareForTest() {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)

	var err error
	ctx.Request, err = http.NewRequest("", "", nil)
	if err != nil {
		panic(err)
	}

	d.Recorder = recorder
	d.Ctx = ctx
	d.Engine = engine
}

/*
Meant to be used with tests.NewGinTester. You can create your
own if you'd like, just ensure the data type that is returned
conforms to GinDataInterface
*/
var InitGinData = func() *GinData {
	data := &GinData{}
	return data
}

////////////
// TESTER //
////////////

/*
Create a new Tester for methods requiring a gin test context. You may use
tests.InitGinData as the initDataFunction for convenience. Otherwise, the custom
data structure provided must conform to the interface tests.GinDataInterface or
else your test will panic. You may compose the tests.GinData into your
data object to get this automatically.

This will automatically attach 2 options to the returned tester.Options:

  - Test Prep at priority tests.DefaultSetupPriority: This casts the data to
    GinDataInterface and calls GinDataInterface.PrepareForTest(). This sets up
    the test recorder, context, and engine for the gin tester.

  - Output Prep at priority 0 (just before the test runs): This casts the
    input at arg index 0 to a gin context.
*/
func NewGinTester[C, M, D any](
	buildMocksFunction func(*testing.T) (C, *M),
	initDataFunction func() *D,
) *Tester[C, M, D] {
	tester := emptyTester[C, M, D]()
	tester.buildMocksFunction = buildMocksFunction
	tester.initDataFunction = initDataFunction

	tester.Options = tester.Options.NewOption(DefaultSetupPriority, func(state *TestState[C, M, D]) {
		ginData := convertToGinDataInterface(state.Data)
		ginData.PrepareForTest()
	})

	tester.Options = tester.Options.SetInput_SC(0, func(state *TestState[C, M, D]) interface{} {
		ginData := convertToGinDataInterface(state.Data)
		return ginData.GetCtx()
	})

	return tester
}

////////////////////////
// SET CONTEXT VALUES //
////////////////////////

/*
Sets a specific key-value pair in the gin context prior to execution.
Has Priority = tests.DefaultInputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputCtx(
	key string, value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		ctx := convertToGinDataInterface(state.Data).GetCtx()
		ctx.Set(key, handleDereference(value))
	})
}

/*
Specify a key and a pointer to a value for a key-value pair in the gin context
prior to execution.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCtx_P(
	key string, valuePointer interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		ctx := convertToGinDataInterface(state.Data).GetCtx()
		ctx.Set(key, removeInterfacePointer(valuePointer))
	})
}

/*
Sets a gin ctx key-value pair based on a provided callback that calculates
that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCtx_C(
	key string, callbackFunction func() interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		ctx := convertToGinDataInterface(state.Data).GetCtx()
		ctx.Set(key, callbackFunction())
	})
}

/*
Sets a gin ctx key-value pair based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCtx_SC(
	key string, callbackFunction func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		ctx := convertToGinDataInterface(state.Data).GetCtx()
		ctx.Set(key, callbackFunction(state))
	})
}

/////////////////////
// METHOD AND URLS //
/////////////////////

// Internal method for handling method and url sets
func applyGin_InputMethodAndURL[C, M, D any](state *TestState[C, M, D], method, url interface{}) {
	methodString, ok := method.(string)
	if !ok {
		panic("could not convert method to a string!")
	}
	urlString, ok := url.(string)
	if !ok {
		panic("could not convert url to a string!")
	}

	var err error
	ctx := convertToGinDataInterface(state.Data).GetCtx()
	ctx.Request.Method = methodString
	ctx.Request.URL, err = urlpkg.Parse(urlString)
	if err != nil {
		panic(err)
	}
}

/*
Set the method and URL for a particular test. If net/url.Parse(url) would cause
an error, this will panic. This will also panic if the method and url are not
string types.
Has Priority = tests.DefaultSetupPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputMethodAndURL(
	method, url interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultSetupPriority, func(state *TestState[C, M, D]) {
		applyGin_InputMethodAndURL(state, handleDereference(method), handleDereference(url))
	})
}

/*
Specify pointers to the method and URL for a particular test. If
net/url.Parse(url) would cause an error, this will panic.
Has Priority = tests.DefaultSetupPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputMethodAndURL_P(
	method, url *string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultSetupPriority, func(state *TestState[C, M, D]) {
		applyGin_InputMethodAndURL(state, *method, *url)
	})
}

/*
Set the method and URL for a particular test based on a provided callback that
calculates those values when this option is reached. If net/url.Parse(url) would
cause an error, this will panic.
Has Priority = tests.DefaultSetupPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputMethodAndURL_C(
	callbackFunction func() (method, url string),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultSetupPriority, func(state *TestState[C, M, D]) {
		method, url := callbackFunction()
		applyGin_InputMethodAndURL(state, method, url)
	})
}

/*
Set the method and URL for a particular test based on a provided callback that
calculates those values based on the value of the TestState when this option is
reached. If net/url.Parse(url) would cause an error, this will panic.
Has Priority = tests.DefaultSetupPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputMethodAndURL_SC(
	callbackFunction func(state *TestState[C, M, D]) (method, url string),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultSetupPriority, func(state *TestState[C, M, D]) {
		method, url := callbackFunction(state)
		applyGin_InputMethodAndURL(state, method, url)
	})
}

//////////////////
// QUERY PARAMS //
//////////////////

// Internal function to set a query param
func applyGin_InputQuery(ginInterface interface{}, key string, value interface{}) {
	valueString, ok := value.(string)
	if !ok {
		panic("could not convert value to a string!")
	}

	ctx := convertToGinDataInterface(ginInterface).GetCtx()
	u := ctx.Request.URL.Query()
	u.Add(key, valueString)
	ctx.Request.URL.RawQuery = u.Encode()
}

/*
Directly specify a key-value pair to be included the request query string.
Has Priority = tests.DefaultInputPriority
value arg Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputQuery(
	key string, value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputQuery(state.Data, key, handleDereference(value))
	})
}

/*
Specify a pointer to a key and value of a key-value pair to be included the
request query string.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputQuery_P(
	key string, value *string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputQuery(state.Data, key, *value)
	})
}

/*
Specify a key-value pair to be included the request query string based on a
provided callback that calculates that value when this option is reached.
*/
func (to *TestOptions[C, M, D]) Gin_InputQuery_C(
	key string, callbackFunction func() string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputQuery(state.Data, key, callbackFunction())
	})
}

/*
Specify a key-value pair to be included the request query string based on a
provided callback that calculates that value based on the value of the TestState
when this option is reached.
*/
func (to *TestOptions[C, M, D]) Gin_InputQuery_SC(
	key string, callbackFunction func(state *TestState[C, M, D]) string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputQuery(state.Data, key, callbackFunction(state))
	})
}

/////////////////////
// SET PATH PARAMS //
/////////////////////

// Internal function to set a query param
func applyGin_InputParam(ginInterface interface{}, key string, value interface{}) {
	valueString, ok := value.(string)
	if !ok {
		panic("could not convert value to a string!")
	}

	ctx := convertToGinDataInterface(ginInterface).GetCtx()
	u := ctx.Request.URL.Query()
	u.Add(key, valueString)
	ctx.Request.URL.RawQuery = u.Encode()
}

/*
Set a specific key-value pair in params.
Has Priority = tests.DefaultInputPriority
value arg Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputParam(
	key string, value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputParam(state.Data, key, handleDereference(value))
	})
}

/*
Set a specific key and a pointer to its value in params
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputParam_P(
	key string, value *string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputParam(state.Data, key, *value)
	})
}

/*
Set a specific key-value pair in params based on a provided callback that
calculates that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputParam_C(
	key string, valueFunction func() string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputParam(state.Data, key, valueFunction())
	})
}

/*
Set a specific key-value pair in params based on a provided callback that
calculates that value based on the value of the TestState when this option is
reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputParam_SC(
	key string, valueFunction func(state *TestState[C, M, D]) string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputParam(state.Data, key, valueFunction(state))
	})
}

////////////////////
// SET BODY VALUE //
////////////////////

// Internal function to write a body value
func applyGin_InputBody[C, M, D any](
	state *TestState[C, M, D],
	value interface{},
) {
	ctx := convertToGinDataInterface(state.Data).GetCtx()

	readCloser, ok := value.(io.ReadCloser)
	if ok {
		ctx.Request.Body = readCloser
	} else {
		body := io.NopCloser(bytes.NewReader(getJsonBytes(value)))
		ctx.Request.Body = body
	}
}

/*
Set the body to the provided value.
Has Priority = tests.DefaultInputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputBody(
	value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputBody(state, handleDereference(value))
	})
}

/*
Set a pointer to the value the body shoudl take during the test.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputBody_P(
	valuePointer interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputBody(state, removeInterfacePointer(valuePointer))
	})
}

/*
Set the body to the provided value based on a provided callback that calculates
that value when this option is reached.
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputBody_C(
	f func() []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputBody(state, f())
	})
}

/*
Set the body to the provided value based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputBody_SC(
	f func(state *TestState[C, M, D]) []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputBody(state, f(state))
	})
}

/////////////////
// GIN HEADERS //
/////////////////

// Internal header write function
func applyGin_InputHeaders[C, M, D any](
	state *TestState[C, M, D],
	headersInterface map[string]interface{},
) {
	ctx := convertToGinDataInterface(state.Data).GetCtx()

	headers := map[string]string{}
	for key, valueInterface := range headersInterface {
		valueString, ok := valueInterface.(string)
		if !ok {
			panic("could not convert header value to a string!")
		}
		headers[key] = valueString
	}

	if ctx.Request.Header == nil {
		ctx.Request.Header = http.Header{}
	}

	for key, value := range headers {
		ctx.Request.Header.Set(key, value)
	}
}

/*
Write a header value.
Has Priority = tests.DefaultInputPriority
value arg Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputHeader(
	key string, value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, map[string]interface{}{key: handleDereference(value)})
	})
}

/*
Write a header value.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputHeader_P(
	key string, value *string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, map[string]interface{}{key: *value})
	})
}

/*
Write a header value based on a provided callback that calculates
that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputHeader_C(
	key string, callbackFunction func() string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, map[string]interface{}{key: callbackFunction()})
	})
}

/*
Write a header value based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputHeader_SC(
	key string, callbackFunction func(state *TestState[C, M, D]) string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, map[string]interface{}{key: callbackFunction(state)})
	})
}

/*
Write header values.
Has Priority = tests.DefaultInputPriority
map value args Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_InputHeaders(
	headers map[string]interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, mapMap(headers, func(value interface{}) interface{} { return handleDereference(value) }))
	})
}

/*
Write header values where the values are all pointers.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputHeaders_P(
	headers map[string]*string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, mapMap(headers, func(value *string) interface{} { return *value }))
	})
}

/*
Write header values based on a provided callback that calculates
that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputHeaders_C(
	key string, valueFunction func() string,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, map[string]interface{}{key: valueFunction()})
	})
}

/*
Write header values based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputHeaders_SC(
	headersFunction func(state *TestState[C, M, D]) map[string]interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputHeaders(state, headersFunction(state))
	})
}

/////////////
// COOKIES //
/////////////

// Internal cookie write function
func applyGin_InputCookies[C, M, D any](
	state *TestState[C, M, D],
	cookies []*http.Cookie,
) {
	ctx := convertToGinDataInterface(state.Data).GetCtx()

	for _, cookie := range cookies {
		ctx.Request.AddCookie(cookie)
	}
}

/*
Add a single cookie to the request.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCookie(
	cookie *http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputCookies(state, []*http.Cookie{cookie})
	})
}

/*
Add a single cookie to the request based on a provided callback that calculates
that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCookie_C(
	callbackFunction func() *http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputCookies(state, []*http.Cookie{callbackFunction()})
	})
}

/*
Add a single cookie to the request based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCookie_SC(
	callbackFunction func(state *TestState[C, M, D]) *http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputCookies(state, []*http.Cookie{callbackFunction(state)})
	})
}

/*
Adds many cookies to the request.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCookies(
	cookies []*http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputCookies(state, cookies)
	})
}

/*
Adds many cookies to the request based on a provided callback that calculates
that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCookies_C(
	cookiesFunction func() []*http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputCookies(state, cookiesFunction())
	})
}

/*
Adds many cookies to the request based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) Gin_InputCookies_SC(
	cookiesFunction func(state *TestState[C, M, D]) []*http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applyGin_InputCookies(state, cookiesFunction(state))
	})
}

/////////////////
// OUTPUT CODE //
/////////////////

/*
Ensures the http code that's written to the recorder matches the provided code
Has Priority = tests.DefaultOutputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_OutputCode(
	code interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		recorder := convertToGinDataInterface(state.Data).GetRecorder()
		state.Assertions.Equal(handleDereference(code), recorder.Code)
	})
}

/*
Ensures the http code that's written to the recorder matches the code that is
pointed to.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputCode_P(
	codePointer *int,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		recorder := convertToGinDataInterface(state.Data).GetRecorder()
		state.Assertions.Equal(*codePointer, recorder.Code)
	})
}

/*
Ensures the http code that's written to the recorder matches based on a provided
callback that calculates that value when this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputCode_C(
	callbackFunction func() int,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		recorder := convertToGinDataInterface(state.Data).GetRecorder()
		state.Assertions.Equal(callbackFunction(), recorder.Code)
	})
}

/*
Ensures the http code that's written to the recorder matches based on a provided
callback that calculates that value based on the value of the TestState when
this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputCode_SC(
	callbackFunction func(state *TestState[C, M, D]) int,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		recorder := convertToGinDataInterface(state.Data).GetRecorder()
		state.Assertions.Equal(callbackFunction(state), recorder.Code)
	})
}

/////////////////
// OUTPUT BODY //
/////////////////

// Internal function for managing uotput body checks
func applyGin_OutputBody[C, M, D any](state *TestState[C, M, D], expectedBody interface{}) {
	// Grab the response bytes
	recorder := convertToGinDataInterface(state.Data).GetRecorder()
	responseBytes, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		panic(err)
	}

	// Assert equality
	state.Assertions.Equal(string(getJsonBytes(expectedBody)), string(responseBytes))
}

/*
Ensures the body of the recorder matches the data passed in.
Has Priority = tests.DefaultOutputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Gin_OutputBody(
	expectedBody interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputBody(state, handleDereference(expectedBody))
	})
}

/*
Ensures the body of the recorder matches the data pointed to.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputBody_P(
	expectedBodyPointer interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputBody(state, removeInterfacePointer(expectedBodyPointer))
	})
}

/*
Ensures the body of the recorder matches the data passed in based on a provided
callback that calculates that value when this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputBody_C(
	callbackFunction func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputBody(state, callbackFunction(state))
	})
}

/*
Ensures the body of the recorder matches the data passed in based on a provided
callback that calculates that value based on the value of the TestState when
this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputBody_SC(
	callbackFunction func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputBody(state, callbackFunction(state))
	})
}

////////////////////
// OUTPUT COOKIES //
////////////////////

// Internal cookie verification function
func applyGin_OutputCookies[C, M, D any](
	state *TestState[C, M, D],
	expectedCookies []*http.Cookie,
) {
	recorder := convertToGinDataInterface(state.Data).GetRecorder()
	actualCookies := recorder.Result().Cookies()

	// Make sure the same number of cookies are returned by both
	state.Assertions.Equal(len(expectedCookies), len(actualCookies))

	for _, expectedCookie := range expectedCookies {

		foundExpectedCookie := false
		for _, actualCookie := range actualCookies {

			// Make sure cookies with the same name are identical
			if expectedCookie.Name == actualCookie.Name {
				state.Assertions.Equal(expectedCookie, actualCookie)
				foundExpectedCookie = true
				break
			}

		}

		// Make sure each expected cookie is actually found
		state.Assertions.Equal(true, foundExpectedCookie)
	}
}

/*
Ensures all cookies that are written to the recorder match all provided cookies.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputCookies(
	expectedCookies []*http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputCookies(state, expectedCookies)
	})
}

/*
Ensures all cookies that are written to the recorder match all provided cookies
based on a provided callback that calculates that value when this option is
reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputCookies_C(
	callbackFunction func() []*http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputCookies(state, callbackFunction())
	})
}

/*
Ensures all cookies that are written to the recorder match all provided cookies
based on a provided callback that calculates that value based on the value of
the TestState when this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Gin_OutputCookies_SC(
	callbackFunction func(state *TestState[C, M, D]) []*http.Cookie,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		applyGin_OutputCookies(state, callbackFunction(state))
	})
}

/*
Create a basic cookie with a name and a value. Shortcut for creating a cookie
like this:

	&http.Cookie{
		Name:  name,
		Value: value,
	}
*/
func SimpleCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
	}
}

/////////////
// HELPERS //
/////////////

// Helper to convert data objects to GinDataInterfaces
func convertToGinDataInterface(data interface{}) GinDataInterface {
	ginData, ok := data.(GinDataInterface)
	if !ok {
		panic("GinTester data type doesn't support GetCtx method")
	}
	return ginData
}

func getJsonBytes(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}
