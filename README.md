# encryption
A package that can encrypt &amp; decrypt any given struct with AES256-GCM based on comparing type and string struct

Initializing a new encryption service:
`
encryption := encryption.NewEncryptionService(yourkey in []byte)
`

Encrypting a struct for a database example:
```
//Returns interface that you can directly store in MongoDB
encryptedStruct, err := encryption.EncryptToInterface(data)
if err != nil {
  //Handle error
  fmt.Println(err)
}
```

Encrypting a struct to further use in your logic:
```
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
