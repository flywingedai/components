package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

///////////////
// TEST DATA //
///////////////

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
Create a new Tester for methods requiring a gin test context. The Data
structure must conform to the interface tests.GinDataInterface or else
your test will panic. You may compose the tests.GinData into your
data object to get this automatically.
*/
func NewGinTester[P, C, M, D any](
	newComponentFunction func(P) C,
	buildMocksFunction func(*testing.T) (P, *M),
	initDataFunction func() *D,
) *Tester[P, C, M, D] {
	tester := &Tester[P, C, M, D]{
		newComponentFunction: newComponentFunction,
		buildMocksFunction:   buildMocksFunction,
		initDataFunction:     initDataFunction,
		Options:              &TestOptions[C, M, D]{},
	}

	tester.Options = tester.Options.SetInput(0, func(state *TestState[C, M, D]) interface{} {
		ginData := convertToGinDataInterface(state.Data)
		ginData.PrepareForTest()
		return ginData.GetCtx()
	})

	return tester
}

///////////
// INPUT //
///////////

/*
Write a specific body value to the test. If you need to write this
value based on the options, use WriteGinBody() instead. You can omit
the args to method and url to leave them empty. For most tests, that
will be sufficient.
*/
func (to *TestOptions[C, M, D]) Gin_WriteBodyValue(
	value interface{},
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinBody(state, value, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write a value to the test request body. You can omit the args
to method and url to leave them empty. For most tests, that
will be sufficient.
*/
func (to *TestOptions[C, M, D]) Gin_WriteBody(
	f func(state *TestState[C, M, D]) []interface{},
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			value := f(state)
			writeGinBody(state, value, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write a header to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_WriteHeaderValue(
	key, value string,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinHeaders(state, map[string]string{key: value}, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write header values to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_WriteHeaderValues(
	headers map[string]string,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinHeaders(state, headers, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write a header to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_WriteHeader(
	key string, valueFunction func(state *TestState[C, M, D]) string,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinHeaders(state, map[string]string{key: valueFunction(state)}, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write header values to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_WriteHeaders(
	headersFunction func(state *TestState[C, M, D]) map[string]string,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinHeaders(state, headersFunction(state), methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write a header to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_AddCookieValue(
	cookie *http.Cookie,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinCookies(state, []*http.Cookie{cookie}, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write header values to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_AddCookieValues(
	cookies []*http.Cookie,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinCookies(state, cookies, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write a header to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_AddCookie(
	cookieFunction func(state *TestState[C, M, D]) *http.Cookie,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinCookies(state, []*http.Cookie{cookieFunction(state)}, methodAndUrl...)
		},
	})
	return testOptions
}

/*
Write header values to the request being made.
*/
func (to *TestOptions[C, M, D]) Gin_AddCookies(
	cookiesFunction func(state *TestState[C, M, D]) []*http.Cookie,
	methodAndUrl ...string,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultInputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			writeGinCookies(state, cookiesFunction(state), methodAndUrl...)
		},
	})
	return testOptions
}

////////////
// OUTPUT //
////////////

/*
Ensures the http code that's written to the recorder matches
*/
func (to *TestOptions[C, M, D]) Gin_ValidateCode(
	code int,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			recorder := convertToGinDataInterface(state.Data).GetRecorder()
			state.Assertions.Equal(recorder.Code, code)
		},
	})
	return testOptions
}

/*
Ensures the body of the recorder matches the data passed in.
*/
func (to *TestOptions[C, M, D]) Gin_ValidateBody(
	expectedBody interface{},
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {

			// Grab the response bytes
			recorder := convertToGinDataInterface(state.Data).GetRecorder()
			responseBytes, err := io.ReadAll(recorder.Result().Body)
			if err != nil {
				panic(err)
			}

			// Assert equality
			state.Assertions.Equal(string(responseBytes), string(getJsonBytes(expectedBody)))
		},
	})
	return testOptions
}

/*
Ensures the cookies that are written to the recorder match.
*/
func (to *TestOptions[C, M, D]) Gin_ValidateCookies(
	expectedCookies []*http.Cookie,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
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
		},
	})
	return testOptions
}

/*
Shorthand to validate any portion of the gin response. As long as the
input to each field is something other than the default value, that
field will be checked
*/
func (to *TestOptions[C, M, D]) Gin_Validate(
	ginValidateOptions GinValidateOptions,
) *TestOptions[C, M, D] {
	testOptions := to.Copy()

	if ginValidateOptions.Code != 0 {
		testOptions = testOptions.Gin_ValidateCode(ginValidateOptions.Code)
	}

	if ginValidateOptions.Body != nil {
		testOptions = testOptions.Gin_ValidateBody(ginValidateOptions.Body)
	}

	if ginValidateOptions.Cookies != nil {
		testOptions = testOptions.Gin_ValidateCookies(ginValidateOptions.Cookies)
	}

	return testOptions
}

type GinValidateOptions struct {
	Code    int
	Body    interface{}
	Cookies []*http.Cookie
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

// Gin body write function
func writeGinBody[C, M, D any](
	state *TestState[C, M, D],
	value interface{},
	methodAndUrl ...string,
) {
	var err error
	ctx := convertToGinDataInterface(state.Data).GetCtx()
	body := io.NopCloser(bytes.NewReader(getJsonBytes(value)))

	ctx.Request.Body = body
	if len(methodAndUrl) == 2 {
		ctx.Request.Method = methodAndUrl[0]
		ctx.Request.URL, err = url.Parse(methodAndUrl[1])
		if err != nil {
			panic(err)
		}
	}
}

// Gin header write function
func writeGinHeaders[C, M, D any](
	state *TestState[C, M, D],
	headers map[string]string,
	methodAndUrl ...string,
) {
	var err error
	ctx := convertToGinDataInterface(state.Data).GetCtx()

	for key, value := range headers {
		ctx.Request.Header.Set(key, value)
	}

	if len(methodAndUrl) == 2 {
		ctx.Request.Method = methodAndUrl[0]
		ctx.Request.URL, err = url.Parse(methodAndUrl[1])
		if err != nil {
			panic(err)
		}
	}
}

// Gin cookie write function
func writeGinCookies[C, M, D any](
	state *TestState[C, M, D],
	cookies []*http.Cookie,
	methodAndUrl ...string,
) {
	var err error
	ctx := convertToGinDataInterface(state.Data).GetCtx()

	for _, cookie := range cookies {
		ctx.Request.AddCookie(cookie)
	}

	if len(methodAndUrl) == 2 {
		ctx.Request.Method = methodAndUrl[0]
		ctx.Request.URL, err = url.Parse(methodAndUrl[1])
		if err != nil {
			panic(err)
		}
	}
}

func SimpleCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
	}
}

func getJsonBytes(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}
