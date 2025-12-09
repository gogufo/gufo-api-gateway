// Copyright 2019-2025 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Business Source License 1.1 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License in the LICENSE file at the root of this repository.
//
// As of the Change Date specified in that file, in accordance with the Business Source
// License, use of this software will be governed by the Apache License, Version 2.0.
//
// THIS SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NON-INFRINGEMENT.
package gufodao

import "strings"

func maskemail(email string) string {
	mailarr := strings.Split(email, "@")
	domain := mailarr[1]
	milbody := mailarr[0]
	domainarr := strings.Split(domain, ".")
	domainbody := domainarr[0]
	mailbodymask := milbody[0:1] + "***"
	domainbodymask := domainbody[0:1] + "***" + domainbody[len(domainbody)-1:]
	maskedemail := mailbodymask + "@" + domainbodymask + "." + domainarr[1]
	return maskedemail
}
