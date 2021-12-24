# Globe-Protocol/encryption

</br>

The encryption package of globe protocol uses the `reflect` package to mainly encrypt structs but also string values. 

- The package can be used by creating a new service of it for easy implementation in different structures
- Initializing the service requires input of your encryption key so that you are not bound to a certain encryption key solution
- Encrypting & Decrypting single string values
- Encrypting & Decrypting single byte values
- Encrypt structs to interface for direct storage to any NoSQL database
- Encrypting structs to a JSON byte interface for further use in logic
- Tag system to disable encryption for certain fields that cannot be encrypted

</br>

</br>

## Table of contents

- [**Installation**]()
- [**Features**]()
  - [Creating an Encryption Service]()
  - [Encryption & Decryption of Strings]()
  - [Encryption & Decryption of Bytes]()
  - [Encryption of structs to JSON]()
  - [Encryption of structs to Interface]()
  - [Decrypting encrypted structs back to Original Struct]()

</br>

</br>

## Installation

Run the following command in any Go project:

`go get github.com/globe-protocol/encryption`

</br>

</br>

## Features

### Creating an Encryption Service

```go
func NewEncryptionService(key []byte) EncryptionService
```

The function allows users to create a new instance of the encryption service using their own key inputted as []byte with 32 bits. This so that it can be easily implemented in any Go project structure.

</br>

#### Example

```go
//Requires 32-bit key input
encryptionService := aes256.NewEncryptionService([]byte("/f532*15=5j145/245n*qw21n9q146/-"))
```

We call the NewEncryption function using a 32-bit long string converted to []byte as encryption key input. This will return a new encryption service implementing all the logic functions of the package.

</br>

</br>

### Encryption & Decryption of  Strings

```go
func (e *encryptionService) EncryptStr(str string) ([]byte, error)
```

The function requires input of any string and will try to encrypt it. If for any reason the encryption function fails it will output `nil` and the occurring error. If it passes it will give the encrypted []byte output and nil for error.

```go
func (e *encryptionService) DecryptStr(b []byte) (string, error)
```

Following up the string encryption function this function takes a byte array. When it fails the output will be `""` and the occurring error. When it passes it will output the formerly encrypted string and `nil` for the error.

</br>

#### Example

We'll use the above created encryptionService to show further examples.

```go
//Output will be long encrypted byte array representing string
encryptedBytes, err := encryptionService.EncryptStr("example string")
if err != nil {
    fmt.Println(err) //Handle error in desired way
}

//Output will "example string" just as the input of the encrypt func
decryptedString, err := encryptionService.DecryptStr(encryptedBytes)
if err != nil {
    fmt.Println(err) //Handle error in desired way
}
```

</br>

</br>

### Encryption & Decryption of Bytes

```go
func (e *encryptionService) EncryptByt(str string) ([]byte, error)
```

The function requires input of any []byte and will try to encrypt it. If for any reason the encryption function fails it will output `nil` and the occurring error. If it passes it will give the encrypted []byte output and nil for error.

```go
func (e *encryptionService) DecryptByt(b []byte) ([]byte, error)
```

Following up the byte encryption function this function takes a byte array. When it fails the output will be `nil` and the occurring error. When it passes it will output the formerly encrypted []byte and `nil` for the error.

</br>

#### Example

We'll use the above created encryptionService to show further examples.

```go
//Notice that string is converted to []byte
encryptedBytes, err := encryptionService.EncryptByt([]byte("example string")
if err != nil {
    fmt.Println(err) //Handle error in desired way
}

//Output will represent "example string" in []byte
decryptedBytes, err := encryptionService.DecryptByt(encryptedBytes)
if err != nil {
    fmt.Println(err) //Handle error in desired way
}
```

</br>

</br>

### Encryption of structs to JSON

```go
func (e *encryptionService) EncryptToJSON(eData interface{}) ([]byte, error)
```

The function requires an input of any **1-dimensional struct**. It currently does not support multi-dimensional structs. If the operation is successful the function will output encrypted []bytes representing the interface of the input and a `nil` for the error output. If it fails it will output `nil` and the occurring error.

</br>

#### Example

```go
type TestStructure {
    Id 			 string	`bson:"_id" encrypted:"false"`
    Value 		 string	`bson:"value"`
    NumberValue	  int	 `bson:"number_value"`
}

//Create test data to encrypt
testData := TestStructure{
    Id: 		"38r4jum241q2981",
    Value:		"value",
    NumberValue: "4269",
}

//Returns JSON bytes representing encrypted interface
encryptedBytes, err := encryption.EncryptToJSON(testData)
if err != nil {
	fmt.Println(err) //Handle error in desired way
}
```

Now that we have encrypted our structure we have encrypted JSON bytes representing our structure. To further use this in our logic we still have to map it back to an interface. Obviously now all the fields that we wanted to be encrypted are []byte. So this means that this will not fit back into our old structure. For this reason we will need a encrypted version of our model. Lets call it EncryptedTestStructure.

```go
type EncryptedTestStructure {
    Id			string	`bson:"_id" encrypted:"false"`
    Value		[]byte 	`bson:"value"`
    NumberValue	 []byte	 `bson:"number_value"`
} // Notice that all fields that are encrypted are now []byte instead of their former type
```

Now that we have created this structure lets use it together with the `json` package to map our bytes back to a structure that we can use in our logic.

```go
var encryptedStruct EncryptedTestStructure

err = json.Unmarshal(encryptedBytes, &encryptedStruct)
if err != nil {
  fmt.Println(err) //Handle error in desired way
}
```

Now we can use the encryptedStruct variable to further use in our logic.

</br>

</br>

### Encryption of structs to Interface

```go
func (e *encryptionService) EncryptToInterface(eData interface{}) (map[string]interface{}, error)
```

The EncryptToInterface function again takes any **1-dimensional struct**. If the operation is successful it will output a `map[string]interface` representing the encrypted version of the input and `nil` for the error output. If it fails it will respond with `nil` and the corresponding error.

</br>

#### Example

```go
//Returns interface that you can directly store in MongoDB
encryptedStruct, err := encryption.EncryptToInterface(structure)
if err != nil {
	fmt.Println(err) //Handle error in desired way
}
```

This encryptedStruct value can now directly be used to save into a MongoDB database.

</br>

</br>

### Decrypting encrypted structs back to Original Struct

```go
func (e *encryptionService) Decrypt(encryptedData interface{}, desiredOutput interface{}) (interface{}, error)
```

The Decrypt function takes in a encrypted structure and an empty structure that should be filled with the decrypted fields. The output will be an interface that you can map back to a structure with the same type as the desiredOutput param and a corresponding error if it fails.

#### Example

Lets say we used the encrypt function in chapter above to encrypt a structure and we saved that to a MongoDB. Now we want to get that structure back from the database and decrypt it to show it to a user. First of I will explain the three structures necessary to create, get and respond data.

```go
//Params is the original struct and is used to map the fields back to their original types
type Params struct {
	Id           string `bson:"_id"`
	Number       int    `bson:"number"`
	Availability bool   `bson:"availabillity"`
	Testvar      string `bson:"testvar"`
}

//This structure is used to directly store into a MongoDB as a map[string]interface{}
type CreateDataParams struct {
	Id           string `bson:"_id" encrypted:"false"`
	Number       string `bson:"number"`
	Availability string `bson:"availabillity"`
	Testvar      string `bson:"testvar"`
}

//This struct is used to get the encrypted data from the database
type GetDataParams struct {
    //Notice that the field kept its type since it is not encrypted all other fields are []byte
    Id           string `bson:"_id"` 	
	Number       []byte `bson:"number"`
	Availability []byte `bson:"availabillity"`
	Testvar      []byte `bson:"testvar"`
}
```

Now that we have the structures setup lets see how we can get the data back from the database.

```go
//Create variable that we will use a "type definition structure"
var typeStruct Params
decryptedInterface, err := encryptionService.Decrypt(encryptedStructure, typeStruct)
if err != nil {
    fmt.Println(err) //Handle error in desired way
}

//Now that we have the decrypted interface lets map it back to a struct that we can use
//We use the decryptedInterface and we try to map it back to a struct with the type Params that we defined as our typeStruct
returnStruct, ok := reflect.Indirect(reflect.ValueOf(decryptedInterface)).Interface().(Params)
if !ok {
    fmt.println("Could not map interface back to struct") //Handle error in desired way
}
```

If all things pass you should now have your decrypted struct back with the same types as that you first encrypted it with.
