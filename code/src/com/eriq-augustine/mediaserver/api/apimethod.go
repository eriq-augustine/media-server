package api;

/*
 * An ApiMethod is a description of a method that can be called on the API.
 * It includes a path, whether or nor to authenticate, handler, and parameters to the handler.
 *
 * ApiMethods make it easy and safe to create API methods.
 * Parameters (existence, required, and type), http semantics, and authorization is all handled for you.
 *
 * Parmaters defined in the Params field of an ApiMethod will be validated and passed into the ApiMethod's
 * handler in the order that they were defined.
 * The name in the definition will be the name of the query http request parameter.
 * This value can come from query parameters or POST form values.
 * All types must match EXACTLY (eg integers must be 'int' not 'int64' or '*int').
 * Parameters my either be ints (API_PARAM_TYPE_INT) or strings (API_PARAM_TYPE_STRING).
 * Feel free to pass JSON in your string.
 * All parameters will be trimmed of whitespace before processing.
 *
 * Note that ApiMethods are not allowed to have empty strings as parameters.
 * An empty string will be treated as a missing parameter.
 * In a similar vein, non-required ints are dangerous because 0 will be used as the empty value.
 * If you need to have a non-required int, consider typing it as a string and inspecting it manually.
 *
 * Because refection in Go does not allow you to find the parameter's name, we must rely on order.
 * In addition to explicitly defined parameters, your handler can have up to four implicit parameters.
 * These parameters may appear in ANY order and you may pick and choose the ones you want (or none).
 *  - userId api.UserId - The id of the user making the request (required authentication).
 *  - token api.Token - The token of the user making the request (required authentication).
 *  - request *http.Request - The http request.
 *  - response http.ResponseWriter - The http response (you should only use this in extreme cases).
 * Remember that in Go, we cannot get parameter names.
 * So you may call these parameters whatever you want, they are made unique by their types.
 * request and response are obvious, but userId and token are a little more unusual.
 * These two are typed to be an int and string respectively, and are only typed special to uniquely identify them.
 * The rest of the codebase expects int and string for their types.
 *
 * The return value for the handler is very flexible.
 * Up to three values can be returned:
 *  - interface{} - Usually a message (see com/eriq-augustine/mediaserver/message).
 *                  This will be turned to JSON and will become the http response.
 *                  Feel free to pass something like "" if you are also passing an error.
 *  - int - An http response code (eg http.StatusOK or http.StatusBadRequest).
 *          If 0, then the code will be inferred from the context.
 *  - error - Any error that occurred.
 *            In the case of an error, the response (interface{}) will be ignored and a failure response will be issued.
 *            The http status will still be honored.
 * Once again, you can specify anywhere between zero and these three return values.
 * The return values must be typed exactly.
 * However, they may be returned in any order.
 *
 * Before being added to any routers, Validate() should be called on all ApiMethods.
 * Validate will check both the method definition, the parameter semantics, and return value semantics.
 * Validate() will panic on any error.
 * Although still not runtime, this should immediatly halt the server if any method is mis-configured.
 */

import (
   "fmt"
   "net/http"
   "reflect"
   "strconv"
   "strings"

   "com/eriq-augustine/mediaserver/auth"
   "com/eriq-augustine/mediaserver/config"
   "com/eriq-augustine/mediaserver/log"
   "com/eriq-augustine/mediaserver/messages"
   "com/eriq-augustine/mediaserver/util"
   "com/eriq-augustine/mediaserver/util/errors"
)

const (
   MULTIPART_PARSE_SIZE = 4 * 1024 * 1024 // Store the first 4M in memory.
)

const (
   API_PARAM_TYPE_INT = iota
   API_PARAM_TYPE_STRING
)

type ApiMethod struct {
   Path string
   Handler interface{}
   Auth bool
   Params []ApiMethodParam
}

type ApiMethodParam struct {
   Name string
   Type int
   Required bool
}

// Will just panic on error.
func (method ApiMethod) Validate() {
   // Check the definitions.
   if (method.Path == "") {
      log.Panic("Bad path for API handler");
   }

   if (method.Handler == nil) {
      log.Panic(fmt.Sprintf("Nil handler for API handler for path: %s", method.Path));
   }

   for _, param := range(method.Params) {
      if (param.Name == "") {
         log.Panic(fmt.Sprintf("Nil name for param for API handler for path: %s", method.Path));
      }

      if (!(param.Type == API_PARAM_TYPE_INT || param.Type == API_PARAM_TYPE_STRING)) {
         log.Panic(fmt.Sprintf("Param (%s) for API handler (%s) has bad type (%d)", param.Name, method.Path, param.Type));
      }
   }

   // Check parameter semantics.

   var handlerType reflect.Type = reflect.TypeOf(method.Handler);

   var numParams int = handlerType.NumIn();
   var additionalParams = 0;

   for i := 0; i < numParams; i++ {
      var paramType reflect.Type = handlerType.In(i);

      if (paramType.String() == "api.Token") {
         additionalParams++;

         if (!method.Auth) {
            log.Panic(fmt.Sprintf("API handler (%s) requested a token without authentication", method.Path));
         }
      } else if (paramType.String() == "api.UserId") {
         additionalParams++;

         if (!method.Auth) {
            log.Panic(fmt.Sprintf("API handler (%s) requested a user id without authentication", method.Path));
         }
      } else if (paramType.String() == "*http.Request") {
         additionalParams++;
      } else if (paramType.String() == "http.ResponseWriter") {
         additionalParams++;
      } else {
         if (!(paramType.String() == "int" || paramType.String() == "string")) {
            log.Panic(fmt.Sprintf("API handler (%s) has an actual parameter with incorrect type (%s) must be string or int", method.Path, paramType.String()));
         }
      }
   }

   if (numParams != len(method.Params) + additionalParams) {
      log.Panic(fmt.Sprintf("API handler (%s) actually expects %d parameters, but is defined to expect %d (%d defined, %d implicit)", method.Path, numParams, len(method.Params) + additionalParams, len(method.Params), additionalParams));
   }

   // Check the return semantics.
   var numReturns int = handlerType.NumOut();

   if (numReturns > 3) {
      log.Panic(fmt.Sprintf("API handler (%s) has too many return values. Got %d. Maximum is 3.", method.Path, numReturns));
   }

   for i := 0; i < numReturns; i++ {
      var returnType reflect.Type = handlerType.Out(i);

      if (!(returnType.String() == "interface {}" || returnType.String() == "int" || returnType.String() == "error")) {
         log.Panic(fmt.Sprintf("API handler (%s) has an bad return type (%s) muct be interface{}, int, or error", method.Path, returnType.String()));
      }
   }
}

func (method ApiMethod) Middleware() func(response http.ResponseWriter, request *http.Request) {
   return func(response http.ResponseWriter, request *http.Request) {
      response.Header().Set("Access-Control-Allow-Origin", "*");
      response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS");
      response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization");
      response.Header().Set("Content-Type", "application/json; charset=UTF-8");

      // Skip preflight checks.
      if (request.Method == "OPTIONS") {
         return;
      }

      if (request.URL != nil) {
         log.Debug(request.URL.String());
      }

      responseObj, httpStatus, err := method.HandleAPIRequest(response, request);

      jsonResponse, _ := util.ToJSON(responseObj);
      sendResponse(jsonResponse, err, httpStatus, response);
   }
}

// This handles the API side of the request.
// None of the boilerplate.
func (method ApiMethod) HandleAPIRequest(response http.ResponseWriter, request *http.Request) (interface{}, int, error) {
   var userId int = -1;
   var ok bool;
   var token string = "";

   if (method.Auth) {
      var responseObject interface{};
      ok, userId, token, responseObject = authRequest(request);
      if (!ok) {
         return responseObject, http.StatusUnauthorized, nil;
      }
   }

   ok, args := createArguments(method, UserId(userId), Token(token), response, request);
   if (!ok) {
      return messages.NewGeneralStatus(false, http.StatusBadRequest), http.StatusBadRequest, nil;
   }

   var handlerValue reflect.Value = reflect.ValueOf(method.Handler);
   returns := handlerValue.Call(args);

   return createReturnValues(method, returns);
}

func createReturnValues(method ApiMethod, returns []reflect.Value) (interface{}, int, error) {
   var responseObj interface{} = nil;
   var httpStatus int = 0;
   var err error = nil;

   // Returns are optional.
   for _, val := range(returns) {
      var returnType reflect.Type = val.Type();

      if (returnType.String() == "interface {}") {
         if (!val.IsNil()) {
            responseObj = val.Interface();
         }
      } else if (returnType.String() == "int") {
         httpStatus = int(val.Int());
      } else if (returnType.String() == "error") {
         if (!val.IsNil()) {
            err = val.Interface().(error);
         }
      } else {
         log.Fatal(fmt.Sprintf("Unkown return type (%s) for API handler for path: %s", returnType.String(), method.Path));
      }
   }

   return responseObj, httpStatus, err;
}

// Get all the parameters setup for invocation.
func createArguments(method ApiMethod, userId UserId, token Token, response http.ResponseWriter, request *http.Request) (bool, []reflect.Value) {
   var handlerType reflect.Type = reflect.TypeOf(method.Handler);
   var numParams int = handlerType.NumIn();

   var apiParamIndex = 0;
   var paramValues []reflect.Value = make([]reflect.Value, numParams);

   for i := 0; i < numParams; i++ {
      var paramType reflect.Type = handlerType.In(i);

      // The user id, token, request, and response get handled specially.
      if (method.Auth && paramType.String() == "api.Token") {
         paramValues[i] = reflect.ValueOf(token);
      } else if (method.Auth && paramType.String() == "api.UserId") {
         paramValues[i] = reflect.ValueOf(userId);
      } else if (paramType.String() == "*http.Request") {
         paramValues[i] = reflect.ValueOf(request);
      } else if (paramType.String() == "http.ResponseWriter") {
         paramValues[i] = reflect.ValueOf(response);
      } else {
         // Normal param, fetch the next api parameter and pass it along.
         ok, val := fetchParam(method.Params[apiParamIndex], request);
         if (!ok) {
            return false, []reflect.Value{};
         }

         paramValues[i] = val;
         apiParamIndex++;
      }
   }

   return true, paramValues;
}

func buildApiUrl(path string) string {
   path = strings.TrimPrefix(path, "/");

   return fmt.Sprintf("/api/v%02d/%s", config.GetIntDefault("apiVersion", 0), path);
}

func fetchParam(param ApiMethodParam, request *http.Request) (bool, reflect.Value) {
   // Only the first call will do anything.
   request.ParseMultipartForm(MULTIPART_PARSE_SIZE);

   var stringValue string = strings.TrimSpace(request.FormValue(param.Name));

   if (param.Required && stringValue == "") {
      log.Warn(fmt.Sprintf("Required parameter not found: %s", param.Name));
      return false, reflect.Value{};
   }

   // If we are looking for string, then we are done.
   if (param.Type == API_PARAM_TYPE_STRING) {
      return true, reflect.ValueOf(stringValue);
   }

   // We must be looking for an int (only ints and strings are allowed).

   // First check for an empty non-required int.
   if (stringValue == "") {
      return true, reflect.ValueOf(0);
   }

   intValue, err := strconv.Atoi(stringValue);
   if (err != nil) {
      log.WarnE(fmt.Sprintf("Unable to convert int parameter (%s) from string: '%s'", param.Name, stringValue), err);
      return false, reflect.ValueOf(0);
   }

   return true, reflect.ValueOf(intValue);
}

// Send a response over |response|.
// All API response (|jsonResponse|) should be valid json.
// On error, |jsonResponse| will be ignored.
// In not supplied, the |httpStatus| will become http.StatusInternalServerError on error and
// http.StatusOK on success.
func sendResponse(jsonResponse string, err error, httpStatus int, response http.ResponseWriter) {
   if (err != nil) {
      log.ErrorE("API Error", err);

      if (httpStatus == 0) {
         httpStatus = http.StatusInternalServerError;
      }

      jsonResponse, _ := util.ToJSON(messages.NewGeneralStatus(false, httpStatus));
      response.WriteHeader(httpStatus);
      fmt.Fprintln(response, jsonResponse);
   } else {
      log.Debug("Successful Response: " + jsonResponse);

      if (httpStatus == 0) {
         httpStatus = http.StatusOK;
      }

      response.WriteHeader(httpStatus)
      fmt.Fprintln(response, jsonResponse);
   }
}

// Tries to authorize a request.
// Returns: success, user id, request token, and response object.
// user id and token will only be populated on success.
// response object will only be populated on error.
func authRequest(request *http.Request) (bool, int, string, interface{}) {
   token, ok := getToken(request);

   if (!ok) {
      return false, 0, "", messages.NewRejectedToken(errors.TokenValidationError{errors.TOKEN_VALIDATION_NO_TOKEN});
   }

   // Check for empty tokens.
   if (strings.TrimSpace(token) == "") {
      return false, 0, "", messages.NewRejectedToken(errors.TokenValidationError{errors.TOKEN_VALIDATION_NO_TOKEN});
   }

   userId, err := auth.ValidateToken(token);
   if (err != nil) {
      validationErr, ok := err.(errors.TokenValidationError);
      if (!ok) {
         // Some other (non-validation) error.
         return false, 0, "", messages.NewGeneralStatus(false, http.StatusInternalServerError);
      }
      return false, 0, "", messages.NewRejectedToken(validationErr);
   }

   return true, userId, token, nil;
}

func notFound() (interface{}, int) {
   return messages.NewGeneralStatus(false, http.StatusNotFound), http.StatusNotFound;
}
