package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type GinDataInterface interface {
	GetCtx() *gin.Context
	GetEngine() *gin.Engine
	GetRecorder() *httptest.ResponseRecorder
}

// Simple GIN context intialization
type GinTestData struct {
	Ctx      *gin.Context
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder
}

func (d *GinTestData) GetCtx() *gin.Context {
	return d.Ctx
}
func (d *GinTestData) GetEngine() *gin.Engine {
	return d.Engine
}
func (d *GinTestData) GetRecorder() *httptest.ResponseRecorder {
	return d.Recorder
}

var InitGinTestData = func() *GinTestData {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)
	return &GinTestData{
		Ctx:      ctx,
		Engine:   engine,
		Recorder: recorder,
	}
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
