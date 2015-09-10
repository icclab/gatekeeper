/*
 * Copyright (c) 2015. Zuercher Hochschule fuer Angewandte Wissenschaften
 *  All Rights Reserved.
 *
 *     Licensed under the Apache License, Version 2.0 (the "License"); you may
 *     not use this file except in compliance with the License. You may obtain
 *     a copy of the License at
 *
 *          http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 *     WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 *     License for the specific language governing permissions and limitations
 *     under the License.
 */

/*
 *     Author: Piyush Harsh,
 *     URL: piyush-harsh.info
 */

package main

func InitMsgs() {
	staticMsgs[0] =
		`[
	{
		"metadata": 
		{
			"source": "T-Nova-AuthZ-Service"
		},
		"info":
		[
			{
				"msg": "Welcome to the T-Nova-AuthZ-Service",
				"purpose": "REST API Usage Guide",
				"disclaimer": "It's not yet final!",
				"notice": "Headers and body formats are not defined here."
			}
		],
		"api":
		[
			{
				"uri": "/",
				"method": "GET",
				"purpose": "REST API Structure and Capability Discovery"
			},
			{
				"uri": "/admin/user/",
				"method": "GET",
				"purpose": "Admin API to get list of all users"
			},
			{
				"uri": "/admin/user/",
				"method": "POST",
				"purpose": "Create a new user"
			},
			{
				"uri": "/admin/user/{user-id}",
				"method": "GET",
				"purpose": "Admin API to get detailed info of a particular user"
			},
			{
				"uri": "/admin/user/{user-id}",
				"method": "PUT",
				"purpose": "Admin API to modify details of a particular user"
			},
			{
				"uri": "/admin/user/{user-id}",
				"method": "DELETE",
				"purpose": "Admin API to delete a particular user"
			},
			{
				"uri": "/token/",
				"method": "POST",
				"purpose": "API to request a new service token"
			},
			{
				"uri": "/token/{token-uuid}",
				"method": "GET",
				"purpose": "Get details of this token, lifetime, user or service id of the creator, etc."
			},
			{
				"uri": "/token/{token-uuid}",
				"method": "DELETE",
				"purpose": "Revoke or essentially delete an existing token."
			},
			{
				"uri": "/token/validate/{token-uuid}",
				"method": "GET",
				"purpose": "Validate a existing token, OK/Not-OK type status response"
			},
			{
				"uri": "/auth/{user-id}",
				"method": "GET",
				"purpose": "Authenticate an existing user. OK/Not-OK type response."
			},
			{
				"uri": "/admin/service/",
				"method": "GET",
				"purpose": "Lists all registered services."
			},
			{
				"uri": "/admin/service/",
				"method": "POST",
				"purpose": "Registers a new service with gatekeeper."
			}
		]
	}
]`
	staticMsgs[1] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "No or corrupt data received."
		}
	]
}`
	staticMsgs[2] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "User already exists."
		}
	]
}`
	staticMsgs[3] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "user created successfully",
			"auth-uri": "/auth/xxx",
			"admin-uri": "/admin/user/yyy",
			"id": "zzz"
		}
	]
}`
	staticMsgs[4] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "list of active users"
		}
	],
	"userlist":
	[
		xxx
	]
}`
	staticMsgs[5] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "Incorrect / Missing Header Attributes."
		}
	]
}`
	staticMsgs[6] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "Incorrect Password."
		}
	]
}`
	staticMsgs[7] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Authentication Successful."
        }
    ],
    "tokenlist":
    {
        "id":
		[
			uuid-xxx
		],
        "valid-until":
		[
			time-yyy
		]
    }
}`
	staticMsgs[8] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Token Details."
        }
    ],
    "token":
    {
        "id":"uuid-xxx",
        "valid-until":"time-yyy"
    }
}`
	staticMsgs[9] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Validation / Authorization Successful."
        }
    ]
}`
	staticMsgs[10] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Validation / Authorization Failed."
        }
    ]
}`
	staticMsgs[11] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Registered Service List."
        }
    ],
    "servicelist":
    {
        "service-key":
		[
			uuid-xxx
		],
        "shortname":
		[
			name-yyy
		]
    }
}`
	staticMsgs[12] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "Service with this shortname already exists."
		}
	]
}`
	staticMsgs[13] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "service registered successfully",
			"service-uri": "/service/xxx",
			"service-key": "yyy"
		}
	]
}`
	staticMsgs[14] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "Account details.",
			"username": "xxx",
			"isadmin": "yyy",
			"capabilitylist": "zzz"
		}
	]
}`
	staticMsgs[15] =
		`
{
	"metadata": 
	{
		"source": "T-Nova-AuthZ-Service"
	},
	"info":
	[
		{
			"msg": "User not found."
		}
	]
}`
	staticMsgs[16] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Update Successful."
        }
    ]
}`
	staticMsgs[17] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Update Failed."
        }
    ]
}`
	staticMsgs[18] =
		`
{
    "metadata":
    {
        "source": "T-Nova-AuthZ-Service"
    },
    "info":
    [
        {
            "msg": "Unauthorized Access."
        }
    ]
}`
}
