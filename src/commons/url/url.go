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

// Package commons/url defines url used by pharos-anchor.
package url

// Base returns the base url as a type of string.
func Base() string { return "/api/v1" }

// Base returns the pharos node url as a type of string.
func PharosNode() string { return "/pharos-node" }

// Base returns the deploy url as a type of string.
func Deploy() string { return "/deploy" }

// Base returns the apps url as a type of string.
func Apps() string { return "/apps" }

// Base returns the start url as a type of string.
func Start() string { return "/start" }

// Base returns the stop url as a type of string.
func Stop() string { return "/stop" }

// Base returns the update url as a type of string.
func Update() string { return "/update" }

// Base returns the nodes url as a type of string.
func Nodes() string { return "/nodes" }

// Base returns the groups url as a type of string.
func Groups() string { return "/groups" }

// Base returns the registries url as a type of string.
func Registries() string { return "/registries" }

// Base returns the management url as a type of string.
func Management() string { return "/management" }

// Base returns the monitoring url as a type of string.
func Monitoring() string { return "/monitoring" }

// Base returns the events url as a type of string.
func Events() string { return "/events" }

// Base returns the create url as a type of string.
func Create() string { return "/create" }

// Base returns the join url as a type of string.
func Join() string { return "/join" }

// Base returns the leave url as a type of string.
func Leave() string { return "/leave" }

// Base returns the register url as a type of string.
func Register() string { return "/register" }

// Base returns the unregister url as a type of string.
func Unregister() string { return "/unregister" }

// Base returns the ping url as a type of string.
func Ping() string { return "/ping" }

// Base returns the resource url as a type of string.
func Resource() string { return "/resource" }

// Base returns the performance url as a type of string.
func Performance() string { return "/performance" }

// Base returns the search url as a type of string.
func Search() string { return "/search" }

// Returning Device url as string.
func Device() string { return "/device" }

// Returning Configuration url as string.
func Configuration() string { return "/configuration" }

//Returning Notification url as string.
func Notification() string { return "/notification" }

//Returning Watch url as string.
func Watch() string { return "/watch" }

// Returning Reboot url as string.
func Reboot() string { return "/reboot" }

// Returning Restore url as string.
func Restore() string { return "/restore" }
