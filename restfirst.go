package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

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
var sessKey string
var headerArrayIp string
// var postBodyStruct Post_Body_Struct
// var tokenResponseStruct TokenResponse
var Tokens []Token

func delete_vol_by_id(w http.ResponseWriter, r *http.Request){
	sessKey = r.Header.Get("x-auth-token")
	headerArrayIp = r.Header.Get("x-array-ip")
	uriString := r.RequestURI
	uriSplit := strings.Split(uriString, "/")
	volUUID := uriSplit[4]
	//reqBody, _ := ioutil.ReadAll(r.Body)
	query := "query=\"uuid EQ " + volUUID + "\""
	fmt.Println("Printing query ", query)
	encodedQuery := url.PathEscape(query)
	fmt.Println("Encoded query is : ", encodedQuery)
	getVolUri := "https://" + headerArrayIp + ":8080/api/v1/volumes?" + encodedQuery
	fmt.Println("Get volume query string uri is ", getVolUri)
	getResp, err := HttpRestyGet(getVolUri)
	if err != nil {
		fmt.Println("Error is : ", err)
	}
	fmt.Println("Printing status of get response ", getResp.Status())
	mapGetVolRequest := make(map[string]interface{})
	err2 := json.Unmarshal(getResp.Body(), &mapGetVolRequest)
	if err2 != nil {
		fmt.Println("Error is : ", err2)
	}
	var volName string
	for k, v := range mapGetVolRequest {
		fmt.Println("\n Key is ", k)
		fmt.Println("\n Value is ", v)
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.String:
			fmt.Println("value of key %s is of time string", k)
		case reflect.Slice:
			fmt.Println("Printing vt ", vt)
			fmt.Println("value of key %s is of type slice", k)
			fmt.Println("Printing the slice values ")
			for index, itemCopy  := range v.([]interface{}){
				fmt.Println("Index is ", index)
				fmt.Println("itemCopy is ", itemCopy)
				itemValueType := reflect.TypeOf(itemCopy)
				fmt.Println("For key, type is %s ", itemValueType)
				switch itemValueType.Kind(){
				case reflect.Map:
					for k1, v1 := range itemCopy.(map[string]interface{}){
						fmt.Println("I am printing Key ", k1)
						fmt.Println("I am printing Value ", v1)
						if k1 == "name" {
							volName = v1.(string)
						}
					}
				}
			}
		case reflect.Map:
			vmap := v.(map[string]interface{})
			for confKey, confVal := range vmap {
				confValType := reflect.TypeOf(confVal)
				switch confValType.Kind() {
				case reflect.String:
					fmt.Println("Inside Map, this is a string type for key %s ", confKey)
				case reflect.Map:
					if confKey == "name" {
						fmt.Println("Printing the name of volume", confVal)
						volName = confVal.(string)
					}
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
	fmt.Println("Printing name of the volume after for loop %s", volName)

	deleteUrl := "https://" + headerArrayIp + ":8080/api/v1/volumes/" + volName
	fmt.Println("DELETE url for volume is ", deleteUrl)

	deleteResp, delErr := HttpRestyDelete(deleteUrl)
	if delErr != nil{
		fmt.Println("Error after delete call is : ", delErr)
	}
	fmt.Println("status of delete rest call is ", deleteResp.Status())


}

func get_vol_by_id(w http.ResponseWriter, r *http.Request){
	fmt.Println("I am in get volume by id call")
	fmt.Println("This REST endpoint is : /containers/v1/volumes/{id} ")
	sessKey = r.Header.Get("x-auth-token")
	headerArrayIp = r.Header.Get("x-array-ip")
	uriString := r.RequestURI
	uriSplit := strings.Split(uriString, "/")

	fmt.Println("Printing len of uriSplit", len(uriSplit))
	volUUID := uriSplit[4]
	fmt.Println("Printing volUUID ", volUUID)
	fmt.Println("Passed uri string is ", uriString)
	fmt.Println("Session key string is ", sessKey)
	fmt.Println("Header Array Ip  string is ", headerArrayIp)
	reqBody, _ := ioutil.ReadAll(r.Body)
	reqB := string(reqBody)
	fmt.Println("Received request for create volume in string form is : %v", reqB)
	query := "query=\"uuid EQ " + volUUID + "\""
	fmt.Println("Printing query ", query)
	encodedQuery := url.PathEscape(query)
	fmt.Println("Encoded query is : ", encodedQuery)
	getVolUri := "https://" + headerArrayIp + ":8080/api/v1/volumes?" + encodedQuery
	fmt.Println("Get volume query string uri is ", getVolUri)

	getResp, err := HttpRestyGet(getVolUri)
	if err != nil {
		fmt.Println("Error is : ", err)
	}
	fmt.Println("Printing get response for volume ", getResp.Body())
	fmt.Println("Printing status of get response ", getResp.Status())

	mapGetVolRequest := make(map[string]interface{})
	//here wright get call for volume information from 3PAR
	err2 := json.Unmarshal(getResp.Body(), &mapGetVolRequest)
	if err2 != nil {
		fmt.Println("Error is : ", err2)
	}
	fmt.Println("Here Response for get vol after unmarshal is :", mapGetVolRequest)
	var volName string
	var volId string
	//var volSize int
	//var volDescription string
	var volBaseSnapshotId  float64
	var volClone bool
	volPublished := false
	var volProvType float64 //{1:"FULL", 2:"TPVV", 3:"SNP", 4:"PEER", 5:"UNKNOWN", 6:"TDVV", 7:"DDS"}
	var volComprState float64 //{1:compression is enabled, 2:compression is disabled, 3:compression is turned off, 4:compression is not available}
	var volUserCpg string
	var volCopyType float64 //{1:"BASE", 2:"PHYSICAL_COPY", 3:"VIRTUAL COPY"}
	var volSizeMiB float64

	resultMapGetVol := make(map[string]interface{})

	for k, v := range mapGetVolRequest {
		fmt.Println("\n Key is ", k)
		fmt.Println("\n Value is ", v)
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.String:
			fmt.Println("value of key %s is of time string", k)
		case reflect.Slice:
			fmt.Println("Printing vt ", vt)
			fmt.Println("value of key %s is of type slice", k)
			fmt.Println("Printing the slice values ")
			for index, itemCopy  := range v.([]interface{}){
				fmt.Println("Index is ", index)
				fmt.Println("itemCopy is ", itemCopy)
				itemValueType := reflect.TypeOf(itemCopy)
				fmt.Println("For key, type is %s ", itemValueType)
				switch itemValueType.Kind(){
				case reflect.Map:

					if k == "members"{
						memberMap := itemCopy.(map[string]interface{})
						volName = memberMap["name"].(string)
						volId = memberMap["uuid"].(string)
						volSizeMiB = memberMap["sizeMiB"].(float64)
						//volDescription = memberMap[""].(string)
						volBaseSnapshotId = memberMap["baseId"].(float64)
						volCopyType = memberMap["copyType"].(float64)
						if volCopyType == 2 {
							volClone = true
						} else {
							volClone = false
						}
						volProvType = memberMap["provisioningType"].(float64)
						volComprState = memberMap["compressionState"].(float64)
						volUserCpg = memberMap["uuid"].(string)


						resultMapGetVol["name"] = volName
						resultMapGetVol["id"] = volId
						resultMapGetVol["size"] = volSizeMiB
						resultMapGetVol["description"] = "No description"
						resultMapGetVol["base_snapshot_id"] = volBaseSnapshotId
						resultMapGetVol["clone"] = volClone
						resultMapGetVol["published"] = volPublished
						confMap := make(map[string]interface{})

						confMap["provisioning"] = volProvType
						confMap["compression"] = volComprState
						confMap["userCpg"] = volUserCpg
						resultMapGetVol["config"] = confMap

						fmt.Println("Printing memberMap ", resultMapGetVol)
					}
					for k1, v1 := range itemCopy.(map[string]interface{}){
						fmt.Println("I am printing Key ", k1)
						fmt.Println("I am printing Value ", v1)
						if k1 == "name" {
							volName = v1.(string)
						}
					}
				}
			}

		case reflect.Map:
			vmap := v.(map[string]interface{})
			for confKey, confVal := range vmap {
				confValType := reflect.TypeOf(confVal)
				switch confValType.Kind() {
				case reflect.String:
					fmt.Println("Inside Map, this is a string type for key %s ", confKey)
				case reflect.Map:
					if confKey == "name" {
						fmt.Println("Printing the name of volume", confVal)
						volName = confVal.(string)
					}
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
	//volName := mapGetVolRequest["name"].(string)
	//volName2 := volName.(string)
	fmt.Println("Printing name of the volume after for loop %s", volName)
	fmt.Fprintln(w, resultMapGetVol)

}

func HttpRestyDelete(URI string) (*resty.Response, error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	headerMap := make(map[string]string)
	headerMap["X-HP3PAR-WSAPI-SessionKey"] = sessKey
	headerMap["Content-type"] = "application/json"
	resp, err := client.R().
		SetHeaders(headerMap).
		//SetBody([]byte(postBody)).
		// SetResult(&AuthSuccess{}).
		Delete(URI)

	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

func HttpRestyGet(URI string) (*resty.Response, error){
	client := resty.New()

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//fmt.Println("Body : ", postBody)
	//fmt.Println("I have got the seesion key: ")
	fmt.Println("Printing the global sessKey value ", sessKey)
	fmt.Println("now the header array ip is ", headerArrayIp)
	headerMap := make(map[string]string)
	//commenting below 2 lines as session key is now global key and will be passed with the header
	//sess_key := session_key.Key
	//headerMap["X-HP3PAR-WSAPI-SessionKey"] = sess_key
	headerMap["X-HP3PAR-WSAPI-SessionKey"] = sessKey
	headerMap["Content-type"] = "application/json"

	resp, err := client.R().
		SetHeaders(headerMap).
		//SetBody([]byte(postBody)).
		// SetResult(&AuthSuccess{}).
		Get(URI)

	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}
func HttpSessionPost(URI string, postBody string) (*resty.Response, error) {

	client := resty.New()

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	fmt.Println("Body : ", postBody)
	fmt.Println("I have got the seesion key: ")

	headerMap := make(map[string]string)
	// sess_key := session_key.Key
	//headerMap["X-HP3PAR-WSAPI-SessionKey"] = sess_key
	headerMap["X-HP3PAR-WSAPI-SessionKey"] = sessKey
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
	sessKey = r.Header.Get("x-auth-token")
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

	//createResponseBody := make(map[string]interface{})
	//json.Unmarshal(response.Body(), createResponseBody)
	//fmt.Println("Response body in map form  for create volume call is ", createResponseBody)
	fmt.Println("Response status is : ", response.Status())
	resBodyCreateVolString := make(map[string]interface{})
	if response.Status() == "201 Created"{
		getVolByNameString := "https://15.212.192.252:8080/api/v1/volumes/"
		volName := mapCreateVolRequest["name"].(string)
		getVolByNameUri := getVolByNameString + volName
		response1, err1 := HttpRestyGet(getVolByNameUri)
		fmt.Println("Printing response status of getvolcall : ", response1.Status())
		if err != nil {
			fmt.Println("Received some error and error is :", err1)
		}
		fmt.Println("Response for get vol is :", response1.Body())
		err2 := json.Unmarshal(response1.Body(), &resBodyCreateVolString)
		if err2 != nil {
			fmt.Println("Error is : ", err)
		}
		fmt.Println("Response for get vol after unmarshal is :", resBodyCreateVolString)
	}
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
	fmt.Fprint(w, string(mapVolCreateResponse))
	fmt.Fprint(w, VolCreateResponseMap)
	//createvolResBody := make(map[string]interface{})
	//json.Unmarshal(response.Body(), createvolResBody)

}


func create_array_session(w http.ResponseWriter, r *http.Request) {

	fmt.Println("This REST endpoint is : /containers/v1/tokens")
	fmt.Println("This function is invoked only when it is a POST call against this end point")
	fmt.Println("Printing the content type from header", r.Header.Get("Content-Type"))
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
	myRouter.HandleFunc("/containers/v1/volumes/{id}", get_vol_by_id).Methods("GET")
	myRouter.HandleFunc("/containers/v1/volumes/{id}", delete_vol_by_id).Methods("DELETE")
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
