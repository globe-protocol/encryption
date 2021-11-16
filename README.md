# encryption
A package that can encrypt &amp; decrypt any given struct with AES256-GCM based on comparing type and string struct

## Initializing a new encryption service:
`
encryption := encryption.NewEncryptionService(yourkey in []byte)
`

## Encrypting a struct for a database example:
```go
//Returns interface that you can directly store in MongoDB
encryptedStruct, err := encryption.EncryptToInterface(data)
if err != nil {
  //Handle error
  fmt.Println(err)
}
```

## Encrypting a struct to further use in your logic:
```go
//Returns json bytes
encryptedBytes, err := encryption.EncryptToJSON(testData)
if err != nil {
  fmt.Println(err)
  return
}

//The fields of the struct that are encrypted should be of []byte type
encryptedStruct := encT{}

err = json.Unmarshal(test, &jsontest)
if err != nil {
  fmt.Println(err)
  return
}
```

## Using the `encrypted` tag and creating structs that are compatible:

First the original struct that we will use to create the encryption structs:
```go
type Data struct {
	Id           string  `json:"_id"
	Number       float64 `json:"number"`
	Availability bool    `json:"availabillity"`
	Testvar      int64   `json:"testvar"`
}
```
Encryption structs (for Encrypting to database storage (MongoDB)):

```go
type CreateDataParams struct {
	Id           string `bson:"_id" encrypted:"false"`
	Number       string `bson:"number"`
	Availability string `bson:"availabillity"`
	Testvar      string `bson:"testvar"`
}

type GetDataParams struct {
	Id           string `bson:"_id"` //id
	Number       []byte `bson:"number"`
	Availability []byte `bson:"availabillity"`
	Testvar      []byte `bson:"testvar"`
}
```
### Note that the bson tags can also be replaced with `json`

The first struct shows 4 fields where 3 should be encrypted. The first one however which is the id should not be encrypted since in that case you wouldn't be able to find the object back in the database. For this we have the `encrypted` tag where you specify `encrypted:"false"` behind the field that should not be encrypted.

Notice that the Get struct has `[]bytes` wherever the fields are encrypted. Where they are not you can just use the type you used in your original struct.
