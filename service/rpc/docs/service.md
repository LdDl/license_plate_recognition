# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [object.proto](#object-proto)
    - [BBox](#license_plate_recognition-BBox)
    - [LPRInfo](#license_plate_recognition-LPRInfo)
    - [LPRRequest](#license_plate_recognition-LPRRequest)
    - [LPRResponse](#license_plate_recognition-LPRResponse)
  
- [service.proto](#service-proto)
    - [Service](#license_plate_recognition-Service)
  
- [Scalar Value Types](#scalar-value-types)



<a name="object-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## object.proto



<a name="license_plate_recognition-BBox"></a>

### BBox
Reference information about detection rectangle


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x_left | [int32](#int32) |  |  |
| y_top | [int32](#int32) |  |  |
| height | [int32](#int32) |  |  |
| width | [int32](#int32) |  |  |






<a name="license_plate_recognition-LPRInfo"></a>

### LPRInfo
Information about single license plate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bbox | [BBox](#license_plate_recognition-BBox) |  | License plate location |
| ocr_bboxes | [BBox](#license_plate_recognition-BBox) | repeated | License plate OCR bounding bboxes. Bounding bboxes are sorted by horizontal line Warning: those coordinates are relative to license plate bounding box, not the parent image! |
| text | [string](#string) |  | License plate text |






<a name="license_plate_recognition-LPRRequest"></a>

### LPRRequest
Essential information to process


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| image | [bytes](#bytes) |  | Bytes representation of image (PNG) |
| bbox | [BBox](#license_plate_recognition-BBox) |  | Optional information about image. Could be usefull if client-side already knows where license plate should be located (due some object detections technique) |






<a name="license_plate_recognition-LPRResponse"></a>

### LPRResponse
Response from server


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| license_plates | [LPRInfo](#license_plate_recognition-LPRInfo) | repeated | Set of found license plates with corresponding information |
| elapsed | [float](#float) |  | Number of seconds has taken to proccess license plate detections and OCR |
| message | [string](#string) |  | Optional message from server |
| warning | [string](#string) |  | Optional warning message from server. If it is not empty you probably should investiage such behavior |
| error | [string](#string) |  | Optional error message from server. If it is not empty you should investiage the error |





 

 

 

 



<a name="service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service.proto


 

 

 


<a name="license_plate_recognition-Service"></a>

### Service
Main service

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ProcessImage | [LPRRequest](#license_plate_recognition-LPRRequest) | [LPRResponse](#license_plate_recognition-LPRResponse) | Unary service: single request == single response |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

