package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/go-resty/resty"
	"github.com/gorilla/mux"
)


type Token struct {
	Array_Ip  string `json:"array_ip"`
	User_Name string `json:"user_name"`
	Password  string `json:"password"`
}

type Session_Key struct {
	Key string `json:"key"`
}

/*
type Post_Body_Struct struct {
	User string `json:"user"`
	Password string `json:"password"`
}
*/

/*
type TokenResponse struct {
	Id string `json:"id"`
	ArrayIp string `json:"array_ip"`
	Username string `json:"username"`
	CreationTime string `json:"creation_time"`
	ExpiryTime string `json:"expiry_time"`
	SessionToken string `json:"session_token"`
}
*/

var session_key Session_Key
var token_struct Token
// var postBodyStruct Post_Body_Struct
// var tokenResponseStruct TokenResponse
var Tokens []Token

func HttpSessionPost(URI string, postBody string) (*resty.Response, error) {

	client := resty.New()

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	fmt.Println("Body : ", postBody)
	fmt.Println("I have got the seesion key: ")

	headerMap := make(map[string]string)
	sess_key := session_key.Key
	headerMap["X-HP3PAR-WSAPI-SessionKey"] = sess_key
	headerMap["Content-type"] = "application/json"
	resp, err := client.R().
		SetHeaders(headerMap).
		SetBody([]byte(postBody)).
		// SetResult(&AuthSuccess{}).
		Post(URI)

	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}

}
func HttpPost(URI string, postBody string) (*resty.Response, error) {

	client := resty.New()

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	fmt.Println("Body : ", postBody)
	fmt.Println("I have got the seesion key: ")
	resp, err := client.R().
		SetHeader("Content-type", "application/json").
		SetBody([]byte(postBody)).
		// SetResult(&AuthSuccess{}).
		Post(URI)

	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	fmt.Println("Endpoint Hit")
}


func getConfigMap(conf interface{})(m map[string]string){
	configVal := conf.(map[string]interface{})

	fmt.Println("Printing the conf value ", conf)
	fmt.Println("type of conf is ", reflect.TypeOf(conf))

	mapString := make(map[string]string)
	for key, value := range configVal{
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}

	return mapString
}


func create_volume(w http.ResponseWriter, r *http.Request){

	/*
	Create Request:
	{
    "data": {
        "name": "my-new-volume",
        "size": "1073741824",
        "description": "my first volume",
        "config": {
            "parameter1": "default",
            "parameter2": false
        }
    }
}

	Create Response:
	{
    "data": {
        "id": "063b5de80e54af7a6b0000000000000000000000d0",
        "name": "my-new-volume",
        "size": 1073741824,
        "description": "my first volume",
        "published": false,
        "base_snapshot_id": "",
        "volume_group_id": "073b5de80e54af7a6b000000000000000000000098",
        "config": {
            "parameter1": "default",
            "parameter2": false
        }
    }
}

	*/
	fmt.Println("I am in create volume call, this is a call for creating a volume")
	fmt.Println("This REST endpoint is : /containers/v1/volumes ")

	reqBody, _ := ioutil.ReadAll(r.Body)
	reqB := string(reqBody)
	fmt.Println("Received request for create volume in string form is : %v", reqB)

	mapCreateVolRequest := make(map[string]interface{})
	err := json.Unmarshal(reqBody, &mapCreateVolRequest)
	if err != nil {
		fmt.Println("Error is : ", err)
	}
	/*
	for k, v := range mapCreateVolRequest {
		fmt.Println("\n Key is ", k)
		fmt.Println("\n Value is ", v)
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.String:
			fmt.Println("value of key %s is of time string", k)
		case reflect.Map:
			vmap := v.(map[string]interface{})
			for confKey, confVal := range vmap {
				confValType := reflect.TypeOf(confVal)
				switch confValType.Kind() {
				case reflect.String:
					fmt.Println("Inside Map, this is a string type for key %s ", confKey)
				case reflect.Bool:
					fmt.Println("Inside Map, this is a boolean type for key %s ", confKey)
				default:
					fmt.Println("Inside Map, this is some other type for key %s ", confKey)
				}
			}
			fmt.Println("value for key %s is of map type", k)
		default:
			fmt.Println("Some other type for key %s ", k)
		}
	}
*/
	fmt.Println("Printing the received create volume request \n", mapCreateVolRequest)

	fmt.Println("printing the request type :", reflect.TypeOf(mapCreateVolRequest))

	fmt.Println("I am using this session key for creating volume : ", session_key.Key)
	configMap := make(map[string]string)

	fmt.Println("Printing configMap value ", configMap)
	pbMapForCreateVol := make(map[string]interface{})
	pbMapForCreateVol["name"] = mapCreateVolRequest["name"]
	//pbMapForCreateVol["cpg"] = mapCreateVolRequest["cpg"]
	sizeInt := mapCreateVolRequest["size"]
	pbMapForCreateVol["sizeMiB"] = sizeInt
	config := mapCreateVolRequest["config"]
	//err1 := json.Unmarshal()

	configMap = getConfigMap(config)
	fmt.Println("Receiveed Config map is ", config)
	fmt.Println("Printing the configMap", config)
	pbMapForCreateVol["cpg"] = configMap["cpg"]
	tpvv_string := configMap["tpvv"]
	var tpvv = true
	if tpvv_string == "true"{
		tpvv = true
	}

	fmt.Println("tpvv is ", tpvv)
	// fmt.Println("compression is ", configMap["compression"])
	arrayIP := configMap["arrayIp"]
	fmt.Println("Printing array ip ", arrayIP)
	pbMapForCreateVol["tpvv"] = tpvv
	fmt.Println("I am printing the type of config ", reflect.TypeOf(config))
	fmt.Println("printing  the name received ", pbMapForCreateVol["name"])
	fmt.Println("I am printing the pbMapForCreateVol: \n ", pbMapForCreateVol)
	//fmt.Println("Type of config resuest is ", reflect.TypeOf(pbMapForCreateVol["cpg"]))
	mapPBCreateVol, _ := json.Marshal(pbMapForCreateVol)
	postBodyCreateVolString := string(mapPBCreateVol)
	fmt.Println("post_body_map_string is ", postBodyCreateVolString)
	//construct create uri for volume
	// https://15.212.196.158:8080/api/v1/volumes
	createVolURI := "https://" + arrayIP + ":8080/api/v1/volumes"
	response, err := HttpSessionPost(createVolURI, postBodyCreateVolString)
	fmt.Println("response for create Vol call is ", response)
	fmt.Println("Response body for create volume call is ", response.Body())
	fmt.Println("Response status is : ", response.Status())
	VolCreateResponseMap := make(map[string]interface{})
	VolCreateResponseMap["name"] = ""
	VolCreateResponseMap["size"] = ""
	VolCreateResponseMap["description"] = ""
	configResp := make(map[string]interface{})
	configResp["tpvv"] = ""
	configResp["cpg"] = ""
	VolCreateResponseMap["config"] = configResp
	mapVolCreateResponse, _ := json.Marshal(VolCreateResponseMap)
	CreateVolResponseString := string(mapVolCreateResponse)
	fmt.Println("CreateVolResponseString is ", CreateVolResponseString)
	//createvolResBody := make(map[string]interface{})
	//json.Unmarshal(response.Body(), createvolResBody)

}


func create_array_session(w http.ResponseWriter, r *http.Request) {

	fmt.Println("This REST endpoint is : /containers/v1/tokens")
	fmt.Println("This function is invoked only when it is a POST call against this end point")
	reqBody, _ := ioutil.ReadAll(r.Body)
	reqB := string(reqBody)
	fmt.Println("Received request in string form is : %v", reqB)
	// fmt.Println("reqBody value is %v", reqBody)

	json.Unmarshal(reqBody, &token_struct)
	user := token_struct.User_Name
	password := token_struct.Password
	arrayIp := token_struct.Array_Ip

	fmt.Println("Array ip is : ", arrayIp)

	postBodyMap := make(map[string]string)
	postBodyMap["user"] = user
	postBodyMap["password"] = password

	mapPBS, _ := json.Marshal(postBodyMap)
	postBodyString := string(mapPBS)
	fmt.Println("post_body_map_string is ", postBodyString)

	postUriForSessionToken := "https://" + arrayIp + ":8080/api/v1/credentials"
	fmt.Println("postUriForSessionToken becomes :", postUriForSessionToken)

	fmt.Println("reqBody ioutil.ReadAll(r.Body) type is \n", reflect.TypeOf(reqBody))


	response, err := HttpPost(postUriForSessionToken, postBodyString)
	/*
	// response, err := HttpPost("https://15.212.192.252:8080/api/v1/credentials", "{\"user\":\"3paradm\",\"password\":\"3pardata\"}")
	// defer response.Body.Close()
	// utils.CheckErr(err)
	// read the response body to a variable
	// bodyBytes, _ := ioutil.ReadAll(response.RawResponse.Body)
	// bodyString := string(bodyBytes)
	//print raw response body for debugging purposes
	// fmt.Println("\n\nSneha and William are helping me ", response.RawResponse.Status, "\n\n")
	// out := acquireBuffer()
	// defer releaseBuffer(out)
	// err := json.Indent(out, r.body, "", "   ")
	// fmt.Println("Hey there ", out.String())
	//println("response and string output ", response.String())
	// res_key, _ := json.Marshal(response.String())
	//fmt.Println("response in string form is ", reflect.TypeOf(response))
*/
	fmt.Println("Response body is", response.Body())
	json.Unmarshal(response.Body(), &session_key)
	// session_token := session_key.Key
	fmt.Println("Printing only session key out: ", session_key.Key)
	fmt.Println("Printing only session key struct out: ", session_key)

	id, err := uuid.NewUUID()
	if err !=nil {
		// handle error
	}
	uuidString := id.String()
	fmt.Printf(uuidString)
	//fmt.Fprintf(w, session_token)
	TokenResponseMap := make(map[string]string)
	TokenResponseMap["id"] = uuidString
	TokenResponseMap["array_ip"] = token_struct.Array_Ip
	TokenResponseMap["username"] = token_struct.User_Name
	TokenResponseMap["creation_time"] = ""
	TokenResponseMap["expiry_time"] = ""
	TokenResponseMap["session_token"] = session_key.Key

	mapTRM, _ := json.Marshal(TokenResponseMap)
	TokenResponseString := string(mapTRM)
	fmt.Println("TokenResponseString is ", TokenResponseString)
	//fmt.Fprintf(w, "response for sessionkey struct", session_key)
	fmt.Fprintln(w, TokenResponseString)
	/*
	//str1 := ""
	//fmt.Sprintf(str1, "%s", response.Body())
	//fmt.Println(" RESPONSE : ", str1)

	// fmt.Println(" Response : ", response.String())
	// fmt.Println(" Error : ", err)
	// fmt.Println("helloooooooo", response)

	// fmt.Fprintf(w, "request received is hello %+v", session_key.Key)
	// var tokenstring Token
	// var responsestring response_key

	// json.Unmarshal(response.Body(), session_key)
	// session_val := session_key.Key
	// fmt.Println("Session val is :", session_val)

	// fmt.Fprintf(w, "hey im printing unmarshaled output %v", reqBody)
	// Tokens = append(Tokens, tokenstring)

	// json.NewEncoder(w).Encode(tokenstring)
	// fmt.Println(json.NewEncoder(w).Encode(response))
*/
}

func get_all_tokens(w http.ResponseWriter, r *http.Request) {
	fmt.Println("listing all token requests")
	fmt.Println("printing the session token key value: ", session_key.Key)
	//fmt.Fprintln(w, Session_Key{})
	//fmt.Fprintf(w, "Printing all tokens ", session_key.Key)
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/containers/v1/tokens", create_array_session).Methods("POST")
	myRouter.HandleFunc("/containers/v1/volumes", create_volume).Methods("POST")
	// myRouter.HandleFunc("/csp/containers/v1/volumes/{id}")
	myRouter.HandleFunc("/alltokens", get_all_tokens)
	// fmt.Println("Received response is %v", string(tokenStringValue))
	err := http.ListenAndServe(":10000", myRouter)
	fmt.Println("Hey there i started REST service")
	if err != nil {
		log.Fatal("Listen and Serve  ERROR", err)
	}
}

func main() {
	handleRequest()
}
