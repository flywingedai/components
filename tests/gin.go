package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
	bufferData := bytes.NewBuffer(getJsonBytes(value))
	if len(methodAndUrl) == 2 {
		ctx.Request, err = http.NewRequest(methodAndUrl[0], methodAndUrl[1], bufferData)
	} else {
		ctx.Request, err = http.NewRequest("", "", bufferData)
	}
	if err != nil {
		panic(err)
	}
}

func getJsonBytes(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}
