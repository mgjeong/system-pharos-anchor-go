/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

// Package commons/errors defines error structs of system-edge-manager.
package errors

// Struct NotFoundURL will be used for return case of error
// which value of unknown or invalid url.
type NotFoundURL struct {
	Message string
}

// Error sets an error message of NotFoundURL.
func (e NotFoundURL) Error() string {
	return "unsupported url: " + e.Message
}

// Struct InvalidMethod will be used for return case of error
// which method of request is not provide.
type InvalidMethod struct {
	Message string
}

// Error sets an error message of InvalidMethod.
func (e InvalidMethod) Error() string {
	return "invalid method: " + e.Message
}

// Struct InvalidParam will be used for return case of error
// which value of unknown or invalid type, range in the parameters.
type InvalidParam struct {
	Message string
}

// Error sets an error message of InvalidParam.
func (e InvalidParam) Error() string {
	return "invalid parameter: " + e.Message
}

// Struct InvalidJSON will be used for return case of error
// which value of malformed json format.
type InvalidJSON struct {
	Message string
}

// Error sets an error message of InvalidJSON.
func (e InvalidJSON) Error() string {
	return "invalid json format: " + e.Message
}

// Struct InvalidJSON will be used for return case of error
// which value of invalid ObjectId.
type InvalidObjectId struct {
	Message string
}

// Error sets an error message of InvalidObjectId.
func (e InvalidObjectId) Error() string {
	return "invalid objectId: " + e.Message
}

// Struct NotFound will be used for return case of error
// which object or target can not found.
type NotFound struct {
	Message string
}

// Error sets an error message of NotFound.
func (e NotFound) Error() string {
	return "not found target: " + e.Message
}

// Struct DBConnectionError will be used for return case of error
// which connection failed with db server.
type DBConnectionError struct {
	Message string
}

// Error sets an error message of DBConnectionError.
func (e DBConnectionError) Error() string {
	return "db connection failed: " + e.Message
}

// Struct DBOperationError will be used for return case of error
// which db operation failed(e.g., insert, update, delete).
type DBOperationError struct {
	Message string
}

// Error sets an error message of DBOperationError.
func (e DBOperationError) Error() string {
	return "db operation failed: " + e.Message
}

// Struct IOError will be used for return case of error
// which IO operaion fail like file operation failed or json marshalling failed.
type IOError struct {
	Message string
}

// Error sets an error message of IOError.
func (e IOError) Error() string {
	return "io error: " + e.Message
}

// Struct InternalServerError will be used for return case of error
// which a generic error, given when an unexpected condition was encountered
// and no more specific message is suitable.
type InternalServerError struct {
	Message string
}

// Error sets an error message of InternalServerError.
func (e InternalServerError) Error() string {
	return "internal server error: " + e.Message
}
